package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	team, err := s.repo.CreateOrUpdateTeam(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *Service) GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := s.repo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if team.Members == nil {
		team.Members = []models.TeamMember{}
	}
	return team, nil
}
