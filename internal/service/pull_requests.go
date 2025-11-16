package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func getTwoRandomIds(ids []string) []string {
	l := len(ids)
	result := []string{}
	if l == 0 {
		return result
	}
	indexes := []int{}
	indexes = append(indexes, rand.Intn(l))
	if l > 1 {
		for {
			if i := rand.Intn(l); i != indexes[0] {
				indexes = append(indexes, i)
				break
			}
		}
	}
	for _, i := range indexes {
		result = append(result, ids[i])
	}
	return result
}

func (s *Service) PullRequestCreate(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	// check PR already exists
	_, err := s.repo.GetPullRequestBase(ctx, pr.PullRequestId)
	if err == nil {
		return nil, models.ErrPRAlreadyExists
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("db: error getting pull request: %w", err)
	}
	// get User and check if not exists
	author, err := s.repo.GetUser(ctx, pr.AuthorId)
	if err == sql.ErrNoRows {
		return nil, models.ErrAuthorNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("db: error checking author: %w", err)
	}

	teamMembers, err := s.repo.GetTeamMembers(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	activeMembersIDs := make([]string, 0)
	for _, member := range teamMembers {
		if member.UserId == author.UserId || !member.IsActive {
			continue
		}
		activeMembersIDs = append(activeMembersIDs, member.UserId)
	}
	pr.AssignedReviewers = getTwoRandomIds(activeMembersIDs)
	createdPR, err := s.repo.CreatePullRequestAndAssignReviewers(ctx, pr)
	if err != nil {
		return nil, err
	}
	return createdPR, nil
}

func (s *Service) PullRequestMerge(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	mergedPR, err := s.repo.MergePullRequest(ctx, pr.PullRequestId)
	if err != nil {
		return nil, err
	}
	return mergedPR, nil
}

func (s *Service) PullRequestReassign(ctx context.Context, prID string, oldUserID string) (*models.PullRequest, string, error) {
	pr, newReviewer, err := s.repo.ReAssignPullRequest(ctx, prID, oldUserID)
	if err != nil {
		return nil, "", err
	}
	return pr, newReviewer, nil
}
