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
