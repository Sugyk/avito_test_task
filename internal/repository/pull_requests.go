package repository

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	getUserQuery := `SELECT id, name, team_name, isActive FROM Users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, getUserQuery, id)
	return &user, err
}

func (r *Repository) GetPullRequestBase(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	checkPRQuery := `SELECT id, title, author_id, status FROM PullRequests WHERE id = $1`
	err := r.db.GetContext(ctx, &pr, checkPRQuery, prID)
	return &pr, err
}

func (r *Repository) CreatePullRequestAndAssignReviewers(ctx context.Context, pullRequest *models.PullRequest) (_ *models.PullRequest, err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				r.logger.Error("db: error while rollback commit", "error", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
				r.logger.Error("db: error while commit", "error", err.Error())
			}
		}
	}()

	createPRQuery := `
	INSERT INTO PullRequests(id, title, author_id, status)
	VALUES ($1, $2, $3, 'OPEN')
	RETURNING id, title, author_id, status
	`
	err = tx.GetContext(ctx,
		pullRequest,
		createPRQuery,
		pullRequest.PullRequestId,
		pullRequest.PullRequestName,
		pullRequest.AuthorId,
	)
	if err != nil {
		return nil, fmt.Errorf("db: error creating pr: %w", err)
	}
	if len(pullRequest.AssignedReviewers) > 0 {
		insertReviewersBuilder := squirrel.Insert("PullRequestsUsers").Columns("pr_id", "user_id")
		for _, id := range pullRequest.AssignedReviewers {
			insertReviewersBuilder = insertReviewersBuilder.Values(pullRequest.PullRequestId, id)
		}
		insertReviewersQuery, args, err := insertReviewersBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("db: error building query: %w", err)
		}
		_, err = tx.ExecContext(ctx, insertReviewersQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("db: error insert reviewers: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("db: commit error: %w", err)
	}
	return pullRequest, nil
}

func (r *Repository) MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	getQuery := `SELECT id, title, author_id, status FROM PullRequests WHERE id = $1`
	err := r.db.GetContext(ctx, &pr, getQuery, prID)
	if err == sql.ErrNoRows {
		return nil, models.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving team: %w", err)
	}

	membersQuery := `
		SELECT user_id 
		FROM PullRequestsUsers 
		WHERE pr_id = $1
	`

	reviewers := make([]string, 0)

	err = r.db.SelectContext(ctx, &reviewers, membersQuery, prID)
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving reviewers: %w", err)
	}
	pr.AssignedReviewers = reviewers

	merged_time := time.Now()

	updateMemberQuery := `
		UPDATE PullRequests
		SET status = 'MERGED', merged_at = $1
		WHERE id = $2
		RETURNING status, merged_at
	`
	err = r.db.GetContext(ctx, &pr, updateMemberQuery, merged_time, prID)
	if err != nil {
		return nil, fmt.Errorf("db: error updating team members: %w", err)
	}

	return &pr, nil
}

func (r *Repository) ReAssignPullRequest(ctx context.Context, prID string, oldUserID string) (_ *models.PullRequest, _ string, err error) {
	var pullRequest models.PullRequest
	var newReviewerID string
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("db: start transaction error: %w", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Error("db: error rollback transaction", "error", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil {
				r.logger.Error("db: error closing transaction", "error", err.Error())
			}
		}
	}()

	checkQuery := `SELECT id, title, author_id, status FROM PullRequests WHERE id = $1`
	err = tx.GetContext(ctx, &pullRequest, checkQuery, prID)
	if err == sql.ErrNoRows {
		return nil, "", models.ErrPRNotFound
	}
	if err != nil {
		return nil, "", fmt.Errorf("db: error retrieving team: %w", err)
	}
	if pullRequest.Status != "OPEN" {
		return nil, "", models.ErrReassigningMergedPR
	}
	var user models.User
	checkQuery = `SELECT id, name, team_name FROM Users WHERE id = $1`
	err = tx.GetContext(ctx, &user, checkQuery, oldUserID)
	if err == sql.ErrNoRows {
		return nil, "", models.ErrUserNotFound
	}
	if err != nil {
		return nil, "", fmt.Errorf("db: error retrieving user: %w", err)
	}

	activeMatesIds := make([]string, 0)
	getActiveTeammatesQuery := `
	SELECT id FROM Users
	WHERE team_name = $1 AND isActive = true AND id != $2 
	`
	err = tx.SelectContext(ctx, &activeMatesIds, getActiveTeammatesQuery, user.TeamName, oldUserID)
	if err == sql.ErrNoRows {
		return nil, "", models.ErrNoActiveCandidates
	}
	if err != nil {
		return nil, "", fmt.Errorf("db: error retrieving teammates: %w", err)
	}
	reviewersIds := make([]string, 0)
	getReviewersQuery := `
	SELECT user_id FROM PullRequestsUsers
	WHERE pr_id = $1
	`
	err = tx.SelectContext(ctx, &reviewersIds, getReviewersQuery, prID)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", fmt.Errorf("db: error retrieving reviewers: %w", err)
	}
	if !slices.Contains(reviewersIds, oldUserID) {
		return nil, "", models.ErrUserNotAssignedToPR
	}
	reviewersIds = slices.DeleteFunc(reviewersIds, func(s string) bool {
		return s == oldUserID
	})
	reviewersSet := make(map[string]struct{})

	for _, reviewer := range reviewersIds {
		reviewersSet[reviewer] = struct{}{}
	}
	for _, mate := range activeMatesIds {
		_, ok := reviewersSet[mate]
		if !ok {
			newReviewerID = mate
			break
		}
	}
	if newReviewerID == "" {
		return nil, "", models.ErrNoActiveCandidates
	}

	insertNewReviewerQuery := `
	INSERT INTO PullRequestsUsers(pr_id, user_id)
	VALUES ($1, $2)
	RETURNING user_id
	`
	var checkNewReviewerID string
	err = tx.GetContext(ctx, &checkNewReviewerID, insertNewReviewerQuery, prID, newReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("db: internal error: error inserting new reviewer: %w", err)
	}
	if checkNewReviewerID != newReviewerID {
		return nil, "", fmt.Errorf("db: internal error: error inserting new reviewer: %w", err)
	}

	deleteOldReviewerQuery := `
	DELETE FROM PullRequestsUsers
	WHERE pr_id = $1 AND user_id = $2 
	RETURNING pr_id, user_id
	`

	checkDeleted := struct {
		CheckDeletedPRId       string `db:"pr_id"`
		CheckDeletedReviewerId string `db:"user_id"`
	}{}

	err = tx.GetContext(ctx, &checkDeleted, deleteOldReviewerQuery, prID, oldUserID)
	if err != nil {
		return nil, "", fmt.Errorf("db: internal error: error deleting old reviewer: %w", err)
	}
	if checkDeleted.CheckDeletedPRId != prID || checkDeleted.CheckDeletedReviewerId != oldUserID {
		return nil, "",
			fmt.Errorf(
				"db: internal error: error deleting old reviewer: deleted pr_id, user_id: %s, %s. Expected: (%s, %s)",
				checkDeleted.CheckDeletedPRId,
				checkDeleted.CheckDeletedReviewerId,
				prID, oldUserID,
			)

	}
	err = tx.Commit()
	if err != nil {
		return nil, "", fmt.Errorf("db: commit error:%w", err)
	}
	pullRequest.AssignedReviewers = append(reviewersIds, newReviewerID)
	return &pullRequest, newReviewerID, nil
}
