package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) GetPullRequestBase(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	checkPRQuery := `SELECT id, title, author_id, status FROM PullRequests WHERE id = $1`
	err := r.db.GetContext(ctx, &pr, checkPRQuery, prID)
	if err == sql.ErrNoRows {
		return &pr, models.ErrPRNotFound
	}
	return &pr, err
}

func (r *Repository) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	ReviewersIds := make([]string, 0)
	getReviewersQuery := `
	SELECT user_id 
	FROM PullRequestsUsers
	WHERE pr_id = $1
	`
	err := r.db.SelectContext(ctx, &ReviewersIds, getReviewersQuery, prID)
	if err == sql.ErrNoRows {
		return []string{}, models.ErrNoReviewers
	}
	if err != nil {
		return []string{}, fmt.Errorf("db: error retrieving reviewers: %w", err)
	}
	return ReviewersIds, nil
}

func (r *Repository) CreatePullRequestAndAssignReviewers(ctx context.Context, pullRequest *models.PullRequest) (_ *models.PullRequest, err error) {
	// preparing transaction
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
	// end preparing transaction

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

func (r *Repository) MergePullRequest(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	merged_time := time.Now()

	mergeQuery := `
		UPDATE PullRequests
		SET status = 'MERGED', merged_at = $1
		WHERE id = $2
		RETURNING status, merged_at
	`
	err := r.db.GetContext(ctx, pr, mergeQuery, merged_time, pr.PullRequestId)
	if err != nil {
		return nil, fmt.Errorf("db: error updating team members: %w", err)
	}

	return pr, nil
}

func (r *Repository) ReAssignPullRequest(ctx context.Context, prID string, oldUser *models.User, newReviewerId string) (_ string, err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("db: start transaction error: %w", err)
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

	insertNewReviewerQuery := `
	INSERT INTO PullRequestsUsers(pr_id, user_id)
	VALUES ($1, $2)
	RETURNING user_id
	`
	var checkNewReviewerID string
	err = tx.GetContext(ctx, &checkNewReviewerID, insertNewReviewerQuery, prID, newReviewerId)
	if err != nil {
		return "", fmt.Errorf("db: internal error: error inserting new reviewer: %w", err)
	}
	if checkNewReviewerID != newReviewerId {
		return "", fmt.Errorf("db: internal error: error inserting new reviewer: %w", err)
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

	err = tx.GetContext(ctx, &checkDeleted, deleteOldReviewerQuery, prID, oldUser.UserId)
	if err != nil {
		return "", fmt.Errorf("db: internal error: error deleting old reviewer: %w", err)
	}
	if checkDeleted.CheckDeletedPRId != prID || checkDeleted.CheckDeletedReviewerId != oldUser.UserId {
		return "",
			fmt.Errorf(
				"db: internal error: error deleting old reviewer: deleted pr_id, user_id: %s, %s. Expected: (%s, %s)",
				checkDeleted.CheckDeletedPRId,
				checkDeleted.CheckDeletedReviewerId,
				prID, oldUser.UserId,
			)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("db: commit error:%w", err)
	}
	return newReviewerId, nil
}
