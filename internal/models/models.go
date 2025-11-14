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

// /team/add:
//
//	post:
//	  tags: [Teams]
//	  summary: Создать команду с участниками (создаёт/обновляет пользователей)
//	  requestBody:
//	    required: true
//	    content:
//	      application/json:
//	        schema:
//	          $ref: '#/components/schemas/Team'
//	        example:
//	          team_name: payments
//	          members:
//	            - user_id: u1
//	              username: Alice
//	              is_active: true
//	            - user_id: u2
//	              username: Bob
//	              is_active: true
//	  responses:
//	    '201':
//	      description: Команда создана
//	      content:
//	        application/json:
//	          schema:
//	            type: object
//	            properties:
//	              team:
//	                $ref: '#/components/schemas/Team'
//	          example:
//	            team:
//	              team_name: backend
//	              members:
//	                - user_id: u1
//	                  username: Alice
//	                  is_active: true
//	                - user_id: u2
//	                  username: Bob
//	                  is_active: true
//	    '400':
//	      description: Команда уже существует
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
//	          example:
//	            error:
//	              code: TEAM_EXISTS
//	              message: team_name already exists
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
