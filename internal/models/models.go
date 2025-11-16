package models

import "fmt"

type TeamMember struct {
	User_id  string `json:"user_id" db:"id"`
	Username string `json:"username" db:"name"`
	IsActive *bool  `json:"is_active" db:"isactive"`
}

func (t *TeamMember) Validate() error {
	if t.User_id == "" {
		return fmt.Errorf("user_id is required")
	}
	if t.Username == "" {
		return fmt.Errorf("username is required")
	}
	if t.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	return nil
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func (t *Team) Validate() error {
	if t.TeamName == "" {
		return fmt.Errorf("name is required")
	}
	if len(t.Members) == 0 {
		return fmt.Errorf("Team must contain at least one member")
	}
	for i, member := range t.Members {
		if err := member.Validate(); err != nil {
			return fmt.Errorf("invalid team member %d: %w", i, err)
		}
	}
	return nil
}

type User struct {
	UserId   string `json:"user_id" db:"id"`
	Username string `json:"username" db:"name"`
	TeamName string `json:"team_name" db:"team_name"`
	IsActive bool   `json:"is_active" db:"isactive"`
}

type TeamAddRequest struct {
	Team Team `json:"team"`
}

func (t *TeamAddRequest) Validate() error {
	if err := t.Team.Validate(); err != nil {
		return fmt.Errorf("invalid team data: %w", err)
	}
	return nil
}

type TeamAddResponse201 struct {
	Team Team `json:"team"`
}

type TeamGetResponse200 struct {
	Team Team `json:"team"`
}

type UsersSetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive *bool  `json:"is_active"`
}

func (u *UsersSetIsActiveRequest) Validate() error {
	if u.UserId == "" {
		return fmt.Errorf("user_id is required")
	}
	if u.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	return nil
}

type UsersSerIsActiveResponse200 struct {
	User User `json:"user"`
}

type PullRequest struct {
	PullRequestId     string   `json:"pull_request_id" db:"id"`
	PullRequestName   string   `json:"pull_request_name" db:"title"`
	AuthorId          string   `json:"author_id" db:"author_id"`
	Status            Status   `json:"status" db:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty" db:"created_at"`
	MergedAt          *string  `json:"mergedAt,omitempty" db:"merged_at"`
}

type PullRequestCreateRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

func (p *PullRequestCreateRequest) Validate() error {
	if p.PullRequestId == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	if p.PullRequestName == "" {
		return fmt.Errorf("pull_request_name is required")
	}
	if p.AuthorId == "" {
		return fmt.Errorf("author_id is required")
	}
	return nil
}
func (p *PullRequestCreateRequest) ToPullRequest() *PullRequest {
	return &PullRequest{
		PullRequestId:   p.PullRequestId,
		PullRequestName: p.PullRequestName,
		AuthorId:        p.AuthorId,
	}
}

type PullRequestCreateResponse201 struct {
	Pr PullRequest `json:"pr"`
}

type PullRequestMergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

func (p *PullRequestMergeRequest) Validate() error {
	if p.PullRequestId == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	return nil
}

func (p *PullRequestMergeRequest) ToPullRequest() *PullRequest {
	return &PullRequest{
		PullRequestId: p.PullRequestId,
	}
}

type PullRequestMergeResponse200 struct {
	Pr PullRequest `json:"pr"`
}

type PullRequestReassignRequest struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

func (p *PullRequestReassignRequest) Validate() error {
	if p.PullRequestId == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	if p.OldReviewerId == "" {
		return fmt.Errorf("old_reviewer_id is required")
	}
	return nil
}

type PullRequestReassignResponse200 struct {
	Pr         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type PullRequestShort struct {
	PullRequestId   string `json:"pull_request_id" db:"id"`
	PullRequestName string `json:"pull_request_name" db:"title"`
	AuthorId        string `json:"author_id" db:"author_id"`
	Status          Status `json:"status" db:"status"`
}

type UsersGetReviewResponse200 struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
