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
		return &models.Team{}, fmt.Errorf("error while creating transaction")
	}
	defer tx.Rollback()

	var teamName string
	checkQuery := `SELECT name FROM Teams WHERE name = $1`
	err = tx.GetContext(ctx, &teamName, checkQuery, team.TeamName)
	if err == nil {
		return &models.Team{}, models.ErrTeamExists
	}
	if err != sql.ErrNoRows {
		return &models.Team{}, err
	}

	insertTeamQuery := `INSERT INTO Teams (name) VALUES ($1) RETURNING name`
	err = tx.GetContext(ctx, &teamName, insertTeamQuery, team.TeamName)
	if err != nil {
		return &models.Team{}, err
	}
	var values []interface{}
	for _, member := range team.Members {
		values = append(values, member.User_id, member.Username, teamName, member.IsActive)
	}
	var scopes = []string{}
	placeholder := "($%d, $%d, $%d, $%d)"
	for i := range team.Members {
		scopes = append(
			scopes,
			fmt.Sprintf(placeholder, i*4+1, i*4+2, i*4+3, i*4+4),
		)
	}

	resultPlaceholder := scopes[0]
	if len(scopes) > 1 {
		resultPlaceholder = strings.Join(scopes, ",")
	}
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
		return &models.Team{}, err
	}
	err = tx.Commit()
	if err != nil {
		return &models.Team{}, err
	}

	return team, nil
}
