package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: error while creating transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			r.logger.Error("db: error while rollback changes", "error", err.Error())
		}
	}()

	var teamName string
	checkQuery := `SELECT name FROM Teams WHERE name = $1`
	err = tx.GetContext(ctx, &teamName, checkQuery, team.TeamName)
	if err == nil {
		return nil, models.ErrTeamExists
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	insertTeamQuery := `INSERT INTO Teams (name) VALUES ($1) RETURNING name`
	err = tx.GetContext(ctx, &teamName, insertTeamQuery, team.TeamName)
	if err != nil {
		return nil, err
	}
	var values []interface{}
	for _, member := range team.Members {
		values = append(values, member.UserId, member.Username, teamName, member.IsActive)
	}
	var scopes = []string{}
	placeholder := "($%d, $%d, $%d, $%d)"
	for i := range team.Members {
		scopes = append(
			scopes,
			fmt.Sprintf(placeholder, i*4+1, i*4+2, i*4+3, i*4+4),
		)
	}

	resultPlaceholder := strings.Join(scopes, ",")
	upsertUserQuery := fmt.Sprintf(
		`
        INSERT INTO Users (id, name, team_name, isActive)
        VALUES %s
        ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		team_name = EXCLUDED.team_name,
		isActive = EXCLUDED.isActive
		`, resultPlaceholder,
	)

	_, err = tx.ExecContext(ctx, upsertUserQuery, values...)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *Repository) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	teamQuery := `SELECT name FROM Teams WHERE name = $1`
	err := r.db.GetContext(ctx, &team.TeamName, teamQuery, teamName)
	if err == sql.ErrNoRows {
		return nil, models.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving team: %w", err)
	}

	membersQuery := `
		SELECT id, name, isActive 
		FROM Users 
		WHERE team_name = $1
	`
	var members []models.TeamMember
	err = r.db.SelectContext(ctx, &members, membersQuery, teamName)
	if err != nil {
		return nil, fmt.Errorf("db: error retrieving team members: %w", err)
	}

	team.Members = members

	return &team, nil
}
