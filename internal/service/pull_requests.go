package service

import (
	"context"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (s *Service) PullRequestCreate(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	createdPR, err := s.repo.CreatePullRequestAndAssignReviewers(ctx, pr)
	if err != nil {
		return nil, err
	}
	return createdPR, nil
}

func (s *Service) PullRequestMerge(ctx context.Context, pr *models.PullRequest) (models.PullRequest, error) {
	// TODO: implement the business logic to merge a pull request
	return models.PullRequest{}, nil
}

func (s *Service) PullRequestReassign(ctx context.Context, prID string, oldUserID string) (models.PullRequest, string, error) {
	// TODO: implement the business logic to reassign a pull request reviewer
	return models.PullRequest{}, "", nil
}
