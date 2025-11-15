//go:generate mockgen -destination=service_mock.go -source=handler.go -package=handlers

package handlers

import (
	"context"
	"log/slog"

	"github.com/Sugyk/avito_test_task/internal/models"
)

type Service interface {
	CreateOrUpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error)
	UsersSetIsActive(ctx context.Context, userID string, isActive bool) (models.User, error)
	PullRequestCreate(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error)
	PullRequestMerge(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error)
	PullRequestReassign(ctx context.Context, prID string, oldUserID string) (models.PullRequest, string, error)
	UsersGetReview(ctx context.Context, userID string) ([]models.PullRequestShort, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
