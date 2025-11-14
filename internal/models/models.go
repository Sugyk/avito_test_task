package models

import "fmt"

type TeamMember struct {
	User_id  string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (t *TeamMember) Validate() error {
	if t.User_id == "" {
		return fmt.Errorf("%w: user_id is required", nil) // TODO: insert error
	}
	if t.Username == "" {
		return fmt.Errorf("%w: username is required", nil) // TODO: insert error
	}
	return nil
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func (t *Team) Validate() error {
	if t.TeamName == "" {
		return fmt.Errorf("%w: name is required", nil) // TODO: insert error
	}
	if len(t.Members) == 0 {
		return fmt.Errorf("%w: name is required", nil) // TODO: insert error
	}
	for _, member := range t.Members {
		if member.Validate() != nil {
			return fmt.Errorf("%w: invalid team member", nil) // TODO: insert error
		}
	}
	return nil
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type TeamAddRequest struct {
	Team Team `json:"team"`
}

func (t *TeamAddRequest) Validate() error {
	if err := t.Team.Validate(); err != nil {
		return fmt.Errorf("%w: invalid team data", nil) // TODO: insert error
	}
	return nil
}

type TeamAddResponse200 struct {
	Team Team `json:"team"`
}

type TeamGetResponse200 struct {
	Team Team `json:"team"`
}

type UsersSetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (u *UsersSetIsActiveRequest) Validate() error {
	if u.UserId == "" {
		return fmt.Errorf("%w: user_id is required", nil) // TODO: insert error
	}
	return nil
}

type UsersSerIsActiveResponse200 struct {
	User User `json:"user"`
}

// /pullRequest/create:
//
//	post:
//	  tags: [PullRequests]
//	  summary: Создать PR и автоматически назначить до 2 ревьюверов из команды автора
//	  requestBody:
//	    required: true
//	    content:
//	      application/json:
//	        schema:
//	          type: object
//	          required: [ pull_request_id, pull_request_name, author_id ]
//	          properties:
//	            pull_request_id: { type: string }
//	            pull_request_name: { type: string }
//	            author_id: { type: string }
//	        example:
//	          pull_request_id: pr-1001
//	          pull_request_name: Add search
//	          author_id: u1
//	  responses:
//	    '201':
//	      description: PR создан
//	      content:
//	        application/json:
//	          schema:
//	            type: object
//	            properties:
//	              pr:
//	                $ref: '#/components/schemas/PullRequest'
//	          example:
//	            pr:
//	              pull_request_id: pr-1001
//	              pull_request_name: Add search
//	              author_id: u1
//	              status: OPEN
//	              assigned_reviewers: [u2, u3]
//	    '404':
//	      description: Автор/команда не найдены
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
//	    '409':
//	      description: PR уже существует
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
//	          example:
//	            error: { code: PR_EXISTS, message: PR id already exists }

type PullRequest struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            Status   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestCreateRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

func (p *PullRequestCreateRequest) Validate() error {
	if p.PullRequestId == "" {
		return fmt.Errorf("%w: pull_request_id is required", nil) // TODO: insert error
	}
	if p.PullRequestName == "" {
		return fmt.Errorf("%w: pull_request_name is required", nil) // TODO: insert error
	}
	if p.AuthorId == "" {
		return fmt.Errorf("%w: author_id is required", nil) // TODO: insert error
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
		return fmt.Errorf("%w: pull_request_id is required", nil) // TODO: insert error
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
		return fmt.Errorf("%w: pull_request_id is required", nil) // TODO: insert error
	}
	if p.OldReviewerId == "" {
		return fmt.Errorf("%w: old_reviewer_id is required", nil) // TODO: insert error
	}
	return nil
}

type PullRequestReassignResponse200 struct {
	Pr         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}

type PullRequestShort struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          Status `json:"status"`
}

type UsersGetReviewResponse200 struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
