package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"slices"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func getActiveUsersIds(teamMembers []models.User) []string {
	activeMembersIDs := make([]string, 0)
	for _, member := range teamMembers {
		if !member.IsActive {
			continue
		}
		activeMembersIDs = append(activeMembersIDs, member.UserId)
	}
	return activeMembersIDs
}

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
	} else if err != models.ErrPRNotFound {
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
	activeTeamMembersIds := getActiveUsersIds(teamMembers)
	activeTeamMembersIds = slices.DeleteFunc(activeTeamMembersIds, func(s string) bool {
		return s == pr.AuthorId
	})
	pr.AssignedReviewers = getTwoRandomIds(activeTeamMembersIds)
	createdPR, err := s.repo.CreatePullRequestAndAssignReviewers(ctx, pr)
	if err != nil {
		return nil, err
	}
	return createdPR, nil
}

func (s *Service) PullRequestMerge(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	pr, err := s.repo.GetPullRequestBase(ctx, pr.PullRequestId)
	if err == models.ErrPRNotFound {
		return nil, err
	}

	pr.AssignedReviewers, err = s.repo.GetPRReviewers(ctx, pr.PullRequestId)
	if err != nil && err != models.ErrNoReviewers {
		return nil, err
	}
	mergedPR, err := s.repo.MergePullRequest(ctx, pr)
	if err != nil {
		return nil, err
	}
	return mergedPR, nil
}

func (s *Service) PullRequestReassign(ctx context.Context, prID string, oldUserID string) (*models.PullRequest, string, error) {
	var pr = &models.PullRequest{PullRequestId: prID}
	// check PR exists
	pr, err := s.repo.GetPullRequestBase(ctx, pr.PullRequestId)
	if err != nil {
		return nil, "", err
	}
	if pr.Status != "OPEN" {
		return nil, "", models.ErrReassigningMergedPR
	}
	// check if old reviewer exists
	user, err := s.repo.GetUser(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}

	// finding new active reviewer
	teamMembers, err := s.repo.GetTeamMembers(ctx, user.TeamName)
	if err != nil {
		return nil, "", models.ErrNoActiveCandidates
	}
	activeMembersIDs := getActiveUsersIds(teamMembers)

	activeMembersIDs = slices.DeleteFunc(activeMembersIDs, func(s string) bool {
		return s == pr.AuthorId || s == oldUserID
	})

	if len(activeMembersIDs) == 0 {
		return nil, "", models.ErrNoActiveCandidates
	}
	reviewersIds, err := s.repo.GetPRReviewers(ctx, prID)
	if err == models.ErrNoReviewers {
		return nil, "", models.ErrNoActiveCandidates
	}
	if err != nil {
		return nil, "", err
	}
	if !slices.Contains(reviewersIds, oldUserID) {
		return nil, "", models.ErrUserNotAssignedToPR
	}
	reviewersIds = slices.DeleteFunc(reviewersIds, func(s string) bool {
		return s == oldUserID
	})

	reviewersSet := make(map[string]struct{})
	var newReviewerID string
	for _, reviewer := range reviewersIds {
		reviewersSet[reviewer] = struct{}{}
	}
	for _, mate := range activeMembersIDs {
		_, ok := reviewersSet[mate]
		if !ok {
			newReviewerID = mate
			break
		}
	}
	if newReviewerID == "" {
		return nil, "", models.ErrNoActiveCandidates
	}

	//

	newReviewer, err := s.repo.ReAssignPullRequest(ctx, prID, user, newReviewerID)
	if err != nil {
		return nil, "", err
	}
	return pr, newReviewer, nil
}
