package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	getUserQuery := `SELECT id, name, team_name, isActive FROM Users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, getUserQuery, id)
	if err == sql.ErrNoRows {
		return &user, models.ErrUserNotFound
	}
	return &user, err
}

func (r *Repository) GetTeamMembers(ctx context.Context, team_name string) ([]models.User, error) {
	var teamMembers []models.User
	getTeamIDsQuery := `
	SELECT id, name, team_name, isActive FROM Users
	WHERE team_name = $1
	`
	err := r.db.SelectContext(ctx, &teamMembers, getTeamIDsQuery, team_name)
	return teamMembers, err
}

func (r *Repository) UsersSetIsActive(ctx context.Context, userID string, isActive bool) error {
	updateMemberQuery := `
		UPDATE Users
		SET isActive = $1
		WHERE id = $2
		RETURNING isActive
	`
	res, err := r.db.ExecContext(ctx, updateMemberQuery, isActive, userID)
	if err != nil {
		return fmt.Errorf("db: error updating team members: %w", err)
	}
	if n, _ := res.RowsAffected(); n != 1 {
		return fmt.Errorf("db: error updating user: expected affected rows: 1. Got: %d", n)
	}

	return nil
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

func (r *Repository) GetActiveTeamMembersIds(ctx context.Context, team_name string, exclude_id string) ([]string, error) {
	var activeTeamMembersIds []string
	getTeamIDsQuery := `
	SELECT id FROM Users
	WHERE team_name = $1 AND isActive = true AND id != $2
	`
	err := r.db.SelectContext(ctx, &activeTeamMembersIds, getTeamIDsQuery, team_name, exclude_id)
	if err != nil {
		return nil, fmt.Errorf("db: error selecting active members: %w", err)
	}

	return activeTeamMembersIds, nil
}
