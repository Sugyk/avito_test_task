package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) UsersSetIsActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	// TODO: implement the business logic to set user's active status
	return models.User{}, nil
}

func (s *Service) UsersGetReview(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	// TODO: implement the business logic to get user's pull requests for review
	return []models.PullRequestShort{}, nil
}
