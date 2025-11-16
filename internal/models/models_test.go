package models

import (
	"testing"
)

func bool_pointer(x bool) *bool {
	return &x
}
func TestTeamMemberValidate(t *testing.T) {
	tests := []struct {
		name    string
		member  TeamMember
		wantErr bool
	}{
		{
			name: "valid member",
			member: TeamMember{
				User_id:  "123",
				Username: "alice",
				IsActive: bool_pointer(true),
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			member: TeamMember{
				User_id:  "",
				Username: "bob",
			},
			wantErr: true,
		},
		{
			name: "missing username",
			member: TeamMember{
				User_id:  "44",
				Username: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestTeamValidate(t *testing.T) {
	tests := []struct {
		name    string
		team    Team
		wantErr bool
	}{
		{
			name: "valid team",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{User_id: "1", Username: "alice", IsActive: bool_pointer(true)},
					{User_id: "2", Username: "bob", IsActive: bool_pointer(true)},
				},
			},
			wantErr: false,
		},
		{
			name: "missing team name",
			team: Team{
				TeamName: "",
				Members: []TeamMember{
					{User_id: "1", Username: "alice"},
				},
			},
			wantErr: true,
		},
		{
			name: "no members",
			team: Team{
				TeamName: "qa",
				Members:  []TeamMember{},
			},
			wantErr: true,
		},
		{
			name: "invalid member",
			team: Team{
				TeamName: "devops",
				Members: []TeamMember{
					{User_id: "", Username: "bob"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestTeamAddRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     TeamAddRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: TeamAddRequest{
				Team: Team{
					TeamName: "frontend",
					Members: []TeamMember{
						{User_id: "1", Username: "alice", IsActive: bool_pointer(true)},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid team",
			req: TeamAddRequest{
				Team: Team{
					TeamName: "",
					Members:  nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestUsersSetIsActiveRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     UsersSetIsActiveRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: UsersSetIsActiveRequest{
				UserId:   "u1",
				IsActive: bool_pointer(true),
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			req: UsersSetIsActiveRequest{
				UserId:   "",
				IsActive: bool_pointer(false),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestPullRequestCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     PullRequestCreateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: PullRequestCreateRequest{
				PullRequestId:   "pr-1",
				PullRequestName: "fix bug",
				AuthorId:        "u1",
			},
			wantErr: false,
		},
		{
			name: "missing pull_request_id",
			req: PullRequestCreateRequest{
				PullRequestId:   "",
				PullRequestName: "xyz",
				AuthorId:        "u1",
			},
			wantErr: true,
		},
		{
			name: "missing pull_request_name",
			req: PullRequestCreateRequest{
				PullRequestId:   "pr-22",
				PullRequestName: "",
				AuthorId:        "u1",
			},
			wantErr: true,
		},
		{
			name: "missing author_id",
			req: PullRequestCreateRequest{
				PullRequestId:   "pr-22",
				PullRequestName: "upgrade",
				AuthorId:        "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestPullRequestCreateRequest_ToPullRequest(t *testing.T) {
	req := PullRequestCreateRequest{
		PullRequestId:   "pr-10",
		PullRequestName: "add search",
		AuthorId:        "u2",
	}

	pr := req.ToPullRequest()

	if pr.PullRequestId != req.PullRequestId {
		t.Errorf("expected PullRequestId %s, got %s", req.PullRequestId, pr.PullRequestId)
	}
	if pr.PullRequestName != req.PullRequestName {
		t.Errorf("expected PullRequestName %s, got %s", req.PullRequestName, pr.PullRequestName)
	}
	if pr.AuthorId != req.AuthorId {
		t.Errorf("expected AuthorId %s, got %s", req.AuthorId, pr.AuthorId)
	}
	if pr.Status != "" {
		t.Errorf("expected empty Status, got %s", pr.Status)
	}
	if len(pr.AssignedReviewers) != 0 {
		t.Errorf("expected no AssignedReviewers, got %v", pr.AssignedReviewers)
	}
	if pr.CreatedAt != nil {
		t.Errorf("expected CreatedAt nil, got %v", pr.CreatedAt)
	}
	if pr.MergedAt != nil {
		t.Errorf("expected MergedAt nil, got %v", pr.MergedAt)
	}
}

func TestPullRequestMergeRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     PullRequestMergeRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: PullRequestMergeRequest{
				PullRequestId: "pr-100",
			},
			wantErr: false,
		},
		{
			name: "missing pull_request_id",
			req: PullRequestMergeRequest{
				PullRequestId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestPullRequestMergeRequest_ToPullRequest(t *testing.T) {
	req := PullRequestMergeRequest{
		PullRequestId: "pr-999",
	}

	pr := req.ToPullRequest()

	if pr.PullRequestId != req.PullRequestId {
		t.Errorf("expected PullRequestId %s, got %s", req.PullRequestId, pr.PullRequestId)
	}
	if pr.PullRequestName != "" {
		t.Errorf("expected empty PullRequestName, got %s", pr.PullRequestName)
	}
	if pr.AuthorId != "" {
		t.Errorf("expected empty AuthorId, got %s", pr.AuthorId)
	}
	if pr.Status != "" {
		t.Errorf("expected empty Status, got %s", pr.Status)
	}
	if len(pr.AssignedReviewers) != 0 {
		t.Errorf("expected no AssignedReviewers, got %v", pr.AssignedReviewers)
	}
	if pr.CreatedAt != nil {
		t.Errorf("expected CreatedAt nil, got %v", pr.CreatedAt)
	}
	if pr.MergedAt != nil {
		t.Errorf("expected MergedAt nil, got %v", pr.MergedAt)
	}
}
