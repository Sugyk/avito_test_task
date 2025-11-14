package service

import (
	"context"
	"log/slog"

	"github.com/Sugyk/avito_test_task/internal/models"
)

type Repository interface {
	CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
}

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}
