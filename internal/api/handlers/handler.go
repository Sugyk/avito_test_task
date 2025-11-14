//go:generate mockgen -destination=service_mock.go -source=handler.go -package=handlers

package handlers

import (
	"log/slog"

	"github.com/Sugyk/avito_test_task/internal/models"
)

type Service interface {
	CreateOrUpdateTeam(team *models.Team) (models.Team, error)
	GetTeamWithMembers(teamName string) (models.Team, error)
	UsersSetIsActive(userID string, isActive bool) (models.User, error)
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
