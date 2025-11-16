package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) UsersSetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.repo.UsersSetIsActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}
	user.IsActive = isActive
	return user, nil
}

func (s *Service) UsersGetReview(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	userPRs, err := s.repo.GetUsersReview(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userPRs, nil
}
