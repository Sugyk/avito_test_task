package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) UsersSetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var user models.User
	getQuery := `SELECT id, name, team_name, isActive FROM Users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, getQuery, userID)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving team: %w", err)
	}

	updateMemberQuery := `
		UPDATE Users
		SET isActive = $1
		WHERE id = $2
		RETURNING isActive
	`
	err = r.db.GetContext(ctx, &user.IsActive, updateMemberQuery, isActive, userID)
	if err != nil {
		return nil, fmt.Errorf("db: error updating team members: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUsersReview(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	var shortPRs = []models.PullRequestShort{}
	checkQuery := `SELECT id FROM Users WHERE id = $1`
	var checkUserID string
	err := r.db.GetContext(ctx, &checkUserID, checkQuery, userID)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving user: %w", err)
	}

	getPrsQuery := `
	SELECT pr_id AS id, title, author_id, status FROM PullRequestsUsers
	LEFT JOIN PullRequests AS pr ON pr_id = pr.id
	WHERE user_id = $1
	`
	err = r.db.SelectContext(ctx, &shortPRs, getPrsQuery, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("db: error selecting PRs of user: %w", err)
	}
	return shortPRs, nil
}
