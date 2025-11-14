package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	team, err := s.repo.CreateOrUpdateTeam(ctx, team)
	if err != nil {
		return &models.Team{}, err
	}
	return team, nil
}

func (s *Service) GetTeamWithMembers(ctx context.Context, teamName string) (models.Team, error) {
	// TODO: implement the business logic to get a team
	return models.Team{}, nil
}
