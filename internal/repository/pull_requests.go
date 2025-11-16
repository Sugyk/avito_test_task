package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Sugyk/avito_test_task/internal/models"
)

func getTwoRandomIds(ids []string) []string {
	l := len(ids)
	result := []string{}
	if l == 0 {
		return result
	}
	indexes := []int{}
	indexes = append(indexes, rand.Intn(l))
	if l > 1 {
		for {
			if i := rand.Intn(l); i != indexes[0] {
				indexes = append(indexes, i)
				break
			}
		}
	}
	for _, i := range indexes {
		result = append(result, ids[i])
	}
	return result
}

func (r *Repository) CreatePullRequestAndAssignReviewers(ctx context.Context, pullRequest *models.PullRequest) (_ *models.PullRequest, err error) {
	var resultPr = &models.PullRequest{}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				r.logger.Error("db: error while rollback commit", "error", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil {
				r.logger.Error("db: error while commit", "error", err.Error())
			}
		}
	}()
	checkPRQuery := `SELECT id FROM PullRequests WHERE id = $1`
	var prID string
	err = tx.GetContext(ctx, &prID, checkPRQuery, pullRequest.PullRequestId)
	if err == nil {
		return nil, models.ErrPRAlreadyExists
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("db: error checking PR: %w", err)
	}

	var checkAuthorQueryRes struct {
		Id        string `db:"id"`
		Team_name string `db:"team_name"`
	}
	checkAuthorQuery := `SELECT id, team_name FROM Users WHERE id = $1`

	err = tx.GetContext(ctx, &checkAuthorQueryRes, checkAuthorQuery, pullRequest.AuthorId)
	if err == sql.ErrNoRows {
		return nil, models.ErrAuthorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error checking author: %w", err)
	}

	var activeTeamMembersIds []string
	getTeamIDsQuery := `
	SELECT id FROM Users
	WHERE team_name = $1 AND isActive = true AND id != $2
	`

	err = tx.SelectContext(ctx, &activeTeamMembersIds, getTeamIDsQuery, checkAuthorQueryRes.Team_name, checkAuthorQueryRes.Id)
	if err != nil {
		return nil, fmt.Errorf("db: error selecting active members: %w", err)
	}

	createPRQuery := `
	INSERT INTO PullRequests(id, title, author_id, status)
	VALUES ($1, $2, $3, 'OPEN')
	RETURNING id, title, author_id, status
	`
	err = tx.GetContext(ctx,
		resultPr,
		createPRQuery,
		pullRequest.PullRequestId,
		pullRequest.PullRequestName,
		checkAuthorQueryRes.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("db: error creating pr: %w", err)
	}

	reviewers_ids := getTwoRandomIds(activeTeamMembersIds)
	if len(reviewers_ids) > 0 {

		insertReviewersBuilder := squirrel.Insert("PullRequestsUsers").Columns("pr_id", "user_id")
		for _, id := range reviewers_ids {
			insertReviewersBuilder = insertReviewersBuilder.Values(pullRequest.PullRequestId, id)
		}

		insertReviewersQuery, args, _ := insertReviewersBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
		_, err = tx.ExecContext(ctx, insertReviewersQuery, args...)
		if err != nil {
			return nil, fmt.Errorf("db: error insert reviewers: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("db: commit error: %w", err)
	}
	resultPr.AssignedReviewers = reviewers_ids
	return resultPr, nil
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
