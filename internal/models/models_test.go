package models

import (
	"testing"
)

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
				IsActive: true,
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
					{User_id: "1", Username: "alice"},
					{User_id: "2", Username: "bob"},
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
						{User_id: "1", Username: "alice"},
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
