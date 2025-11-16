package service

import (
	"context"
	"log/slog"

	"github.com/Sugyk/avito_test_task/internal/models"
)

type Repository interface {
	CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	UsersSetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	CreatePullRequestAndAssignReviewers(ctx context.Context, pullRequest *models.PullRequest) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error)
	ReAssignPullRequest(ctx context.Context, prID string, oldUserID string) (*models.PullRequest, string, error)
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
