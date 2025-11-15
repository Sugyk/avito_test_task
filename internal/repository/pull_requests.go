package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

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
			tx.Rollback()
		} else {
			tx.Commit()
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
	insertReviewersBuilder := squirrel.Insert("PullRequestsUsers").Columns("pr_id", "user_id")
	for _, id := range reviewers_ids {
		insertReviewersBuilder.Values(prID, id)
	}

	insertReviewersQuery, args, err := insertReviewersBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	_, err = tx.ExecContext(ctx, insertReviewersQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("db: error insert reviewers: %w", err)
	}
	resultPr.AssignedReviewers = reviewers_ids
	return resultPr, nil
}
