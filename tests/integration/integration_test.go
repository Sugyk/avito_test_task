package integration

import (
	"net/http"
	"testing"

	"github.com/Sugyk/avito_test_task/internal/models"
	"github.com/stretchr/testify/assert"
)

func bool_pointer(x bool) *bool {
	return &x
}

func TestTeamAdd(t *testing.T) {
	req := models.Team{
		TeamName: "test1",
		Members: []models.TeamMember{
			{
				UserId:   "1",
				Username: "Test1",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "2",
				Username: "Test2",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "3",
				Username: "Test3",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "4",
				Username: "Test4",
				IsActive: bool_pointer(true),
			},
		},
	}

	resp, body := DoPOST(t, "/team/add", req, nil)
	defer resp.Body.Close()

	AssertStatusCode(t, resp, http.StatusCreated)

	var expectedResp = models.TeamAddResponse201{
		Team: models.Team{
			TeamName: "test1",
			Members: []models.TeamMember{
				{
					UserId:   "1",
					Username: "Test1",
					IsActive: bool_pointer(true),
				},
				{
					UserId:   "2",
					Username: "Test2",
					IsActive: bool_pointer(true),
				},
				{
					UserId:   "3",
					Username: "Test3",
					IsActive: bool_pointer(true),
				},
				{
					UserId:   "4",
					Username: "Test4",
					IsActive: bool_pointer(true),
				},
			},
		},
	}
	var createResp models.TeamAddResponse201
	UnmarshalJSON(t, body, &createResp)
	AssertStatusCode(t, resp, 201)
	assert.Equal(t, createResp.Team.TeamName, expectedResp.Team.TeamName)
	for i := range createResp.Team.Members {
		assert.Equal(t, createResp.Team.Members[i].UserId, expectedResp.Team.Members[i].UserId)
		assert.Equal(t, createResp.Team.Members[i].Username, expectedResp.Team.Members[i].Username)
		assert.Equal(t, createResp.Team.Members[i].IsActive, expectedResp.Team.Members[i].IsActive)
	}
}

func TestTeamGetNotFound(t *testing.T) {
	resp, _ := DoGET(t, "/team/get?team_name=unexisting", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
}

func TestTeamGet(t *testing.T) {
	req := models.Team{
		TeamName: "test2",
		Members: []models.TeamMember{
			{
				UserId:   "1",
				Username: "Test1",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "2",
				Username: "Test2",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "3",
				Username: "Test3",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "4",
				Username: "Test4",
				IsActive: bool_pointer(true),
			},
		},
	}
	expectedResp := models.Team{}
	resp, _ := DoPOST(t, "/team/add", req, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusCreated)

	resp, body := DoGET(t, "/team/get?team_name=test2", nil)
	AssertStatusCode(t, resp, http.StatusOK)
	UnmarshalJSON(t, body, &expectedResp)
	AssertStatusCode(t, resp, http.StatusOK)
	assert.Equal(t, req.TeamName, expectedResp.TeamName)
	for i := range req.Members {
		assert.Equal(t, req.Members[i].UserId, expectedResp.Members[i].UserId)
		assert.Equal(t, req.Members[i].Username, expectedResp.Members[i].Username)
		assert.Equal(t, req.Members[i].IsActive, expectedResp.Members[i].IsActive)
	}
}

func TestPullRequestCreate(t *testing.T) {

}
