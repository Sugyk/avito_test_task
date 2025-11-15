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
