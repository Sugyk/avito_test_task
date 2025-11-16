package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Sugyk/avito_test_task/internal/models"
)

func (r *Repository) GetTeamBase(ctx context.Context, team *models.Team) (*models.Team, error) {
	checkQuery := `SELECT name FROM Teams WHERE name = $1`
	err := r.db.GetContext(ctx, &team.TeamName, checkQuery, team.TeamName)
	if err == sql.ErrNoRows {
		return team, models.ErrUserNotFound
	}
	return team, err
}

func (r *Repository) CreateOrUpdateTeam(ctx context.Context, team *models.Team) (_ *models.Team, err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db: error while creating transaction")
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				r.logger.Error("db: error while rollback changes", "error", err.Error())
			}
		} else {
			if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
				r.logger.Error("db: error while rollback changes", "error", err.Error())
			}
		}
	}()

	insertTeamQuery := `INSERT INTO Teams (name) VALUES ($1) RETURNING name`
	res, err := tx.ExecContext(ctx, insertTeamQuery, team.TeamName)
	if err != nil {
		return nil, err
	}
	rowsN, _ := res.RowsAffected()
	if rowsN != 1 {
		return nil, fmt.Errorf("db: inserting team error: affected rows expected: 1, got: %d", rowsN)
	}

	insertQuery := squirrel.Insert("Users").Columns("id", "name", "team_name", "isActive")
	if len(team.Members) > 0 {
		for _, member := range team.Members {
			insertQuery = insertQuery.Values(member.UserId, member.Username, team.TeamName, member.IsActive)
		}
	}
	insertQuery = insertQuery.Suffix(
		`
		ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		team_name = EXCLUDED.team_name,
		isActive = EXCLUDED.isActive
		`,
	)
	insertTeamQuery, args, err := insertQuery.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("db: error building query: %w", err)
	}
	res, err = tx.ExecContext(ctx, insertTeamQuery, args...)
	if err != nil {
		return nil, err
	}
	rowsN, _ = res.RowsAffected()
	if int(rowsN) != len(team.Members) {
		return nil, fmt.Errorf("db: inserting team error: affected rows expected: %d, got: %d", len(team.Members), rowsN)
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
