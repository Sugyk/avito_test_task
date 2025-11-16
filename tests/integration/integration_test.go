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

func addTeam(t *testing.T, team_name string, user_id string, username string, is_active bool) {
	req := models.Team{
		TeamName: team_name,
		Members: []models.TeamMember{
			{
				UserId:   user_id,
				Username: username,
				IsActive: bool_pointer(is_active),
			},
		},
	}
	resp, _ := DoPOST(t, "/team/add", req, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusCreated)

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
		assert.Equal(t, expectedResp.Team.Members[i].UserId, createResp.Team.Members[i].UserId)
		assert.Equal(t, expectedResp.Team.Members[i].Username, createResp.Team.Members[i].Username)
		assert.Equal(t, expectedResp.Team.Members[i].IsActive, createResp.Team.Members[i].IsActive)
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
	assert.Equal(t, expectedResp.TeamName, req.TeamName)
	for i := range req.Members {
		assert.Equal(t, expectedResp.Members[i].UserId, req.Members[i].UserId)
		assert.Equal(t, expectedResp.Members[i].Username, req.Members[i].Username)
		assert.Equal(t, expectedResp.Members[i].IsActive, req.Members[i].IsActive)
	}
}

func TestUsersSetIsActive(t *testing.T) {
	req := models.UsersSetIsActiveRequest{
		UserId:   "1",
		IsActive: bool_pointer(false),
	}
	expectedResp := models.UsersSerIsActiveResponse200{}
	resp, body := DoPOST(t, "/users/setIsActive", req, nil)
	AssertStatusCode(t, resp, http.StatusOK)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.User.UserId, req.UserId)
	assert.Equal(t, expectedResp.User.UserId, req.UserId)
}

func TestPullRequestCreate(t *testing.T) {
	addReq := models.Team{
		TeamName: "test3",
		Members: []models.TeamMember{
			{
				UserId:   "pr-1author",
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
	resp, _ := DoPOST(t, "/team/add", addReq, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusCreated)

	expectedResp := models.PullRequestCreateResponse201{}
	req := models.PullRequestCreateRequest{
		PullRequestId:   "pr-1001",
		PullRequestName: "testpr1",
		AuthorId:        "pr-1author",
	}
	resp, body := DoPOST(t, "/pullRequest/create", req, nil)
	AssertStatusCode(t, resp, http.StatusCreated)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.Pr.PullRequestId, req.PullRequestId)
	assert.Equal(t, expectedResp.Pr.PullRequestId, req.PullRequestId)
}

func TestPullRequestCreateAuthorNotFound(t *testing.T) {
	expectedResp := models.ErrorResponse{}
	req := models.PullRequestCreateRequest{
		PullRequestId:   "pr-1002",
		PullRequestName: "testpr1",
		AuthorId:        "unexisting",
	}
	resp, body := DoPOST(t, "/pullRequest/create", req, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusNotFound)
	UnmarshalJSON(t, body, &expectedResp)

	assert.Equal(t, expectedResp.Error.Code, models.NotFoundErrorCode)
	assert.Equal(t, expectedResp.Error.Message, models.ErrAuthorNotFound.Error())
}

func TestPullRequestCreatePrExists(t *testing.T) {
	expectedResp := models.ErrorResponse{}
	req := models.PullRequestCreateRequest{
		PullRequestId:   "pr-1003",
		PullRequestName: "testpr1",
		AuthorId:        "pr-1author",
	}
	resp, _ := DoPOST(t, "/pullRequest/create", req, nil)
	defer resp.Body.Close()

	AssertStatusCode(t, resp, http.StatusCreated)
	resp, body := DoPOST(t, "/pullRequest/create", req, nil)
	AssertStatusCode(t, resp, http.StatusConflict)
	UnmarshalJSON(t, body, &expectedResp)

	assert.Equal(t, expectedResp.Error.Code, models.PrExistsErrorCode)
	assert.Equal(t, expectedResp.Error.Message, models.ErrPRAlreadyExists.Error())
}

func TestPullRequestMerge(t *testing.T) {
	addTeam(t, "merge-team", "merge-author", "MergeAuthor", true)
	prReq := models.PullRequestCreateRequest{
		PullRequestId:   "pr-merge-1",
		PullRequestName: "MergeTest",
		AuthorId:        "merge-author",
	}
	resp, _ := DoPOST(t, "/pullRequest/create", prReq, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusCreated)

	mergeReq := models.PullRequestMergeRequest{PullRequestId: "pr-merge-1"}
	expectedResp := models.PullRequestMergeResponse200{}
	resp, body := DoPOST(t, "/pullRequest/merge", mergeReq, nil)
	AssertStatusCode(t, resp, http.StatusOK)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.Pr.PullRequestId, mergeReq.PullRequestId)
	assert.Equal(t, expectedResp.Pr.Status, models.StatusMerged)
	assert.NotNil(t, expectedResp.Pr.MergedAt)
}

func TestPullRequestMergeNotFound(t *testing.T) {
	mergeReq := models.PullRequestMergeRequest{PullRequestId: "pr-notfound"}
	expectedResp := models.ErrorResponse{}
	resp, body := DoPOST(t, "/pullRequest/merge", mergeReq, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusNotFound)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.Error.Code, models.NotFoundErrorCode)
	assert.Equal(t, expectedResp.Error.Message, models.ErrPRNotFound.Error())
}

func TestPullRequestReassignNotFound(t *testing.T) {
	reassignReq := models.PullRequestReassignRequest{
		PullRequestId: "pr-nonexist",
		OldReviewerId: "u1",
	}
	expectedResp := models.ErrorResponse{}
	resp, body := DoPOST(t, "/pullRequest/reassign", reassignReq, nil)
	defer resp.Body.Close()
	AssertStatusCode(t, resp, http.StatusNotFound)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.Error.Code, models.NotFoundErrorCode)
}

func TestUsersGetReview(t *testing.T) {
	req := models.Team{
		TeamName: "TestUsersGetReviewtest1",
		Members: []models.TeamMember{
			{
				UserId:   "TestUsersGetReview1",
				Username: "Test1",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "TestUsersGetReview2",
				Username: "Test2",
				IsActive: bool_pointer(true),
			},
		},
	}

	resp, _ := DoPOST(t, "/team/add", req, nil)
	defer resp.Body.Close()

	AssertStatusCode(t, resp, http.StatusCreated)

	prReq := models.PullRequestCreateRequest{
		PullRequestId:   "TestUsersGetReview1",
		PullRequestName: "ReviewTest",
		AuthorId:        "TestUsersGetReview1",
	}
	resp, _ = DoPOST(t, "/pullRequest/create", prReq, nil)
	AssertStatusCode(t, resp, http.StatusCreated)

	prReq1 := models.PullRequestCreateRequest{
		PullRequestId:   "TestUsersGetReview2",
		PullRequestName: "ReviewTest",
		AuthorId:        "TestUsersGetReview1",
	}
	resp, _ = DoPOST(t, "/pullRequest/create", prReq1, nil)
	AssertStatusCode(t, resp, http.StatusCreated)

	expectedResp := models.UsersGetReviewResponse200{}
	resp, body := DoGET(t, "/users/getReview?user_id=TestUsersGetReview2", nil)
	AssertStatusCode(t, resp, http.StatusOK)
	UnmarshalJSON(t, body, &expectedResp)
	assert.Equal(t, expectedResp.UserId, "TestUsersGetReview2")
	assert.Equal(t, len(expectedResp.PullRequests), 2)
	assert.Equal(t, expectedResp.PullRequests[0].PullRequestId, "TestUsersGetReview1")
	assert.Equal(t, expectedResp.PullRequests[1].PullRequestId, "TestUsersGetReview2")
}
