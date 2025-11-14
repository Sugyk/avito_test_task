package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sugyk/avito_test_task/internal/models"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTeamAdd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	teamInput := models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{User_id: "1", Username: "alice"},
		},
	}

	mockService.
		EXPECT().
		CreateOrUpdateTeam(&teamInput).
		Return(teamInput, nil)

	h := NewHandler(mockService, nil)

	body, _ := json.Marshal(models.TeamAddRequest{Team: teamInput})
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TeamAdd(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var out models.TeamAddResponse200
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	require.Equal(t, teamInput.TeamName, out.Team.TeamName)
	require.Len(t, out.Team.Members, 1)
}

func TestTeamAdd_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.TeamAdd(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTeamAdd_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := NewHandler(nil, slog.Default())

	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty team name",
			body: `{"team": {"team_name": "", "members": [{"user_id":"1","username":"bob"}]}}`,
		},
		{
			name: "no members",
			body: `{"team": {"team_name": "backend", "members": []}}`,
		},
		{
			name: "invalid member fields",
			body: `{"team": {"team_name": "backend", "members": [{"user_id":"","username":""}]}}`,
		},
		{
			name: "team object exists but empty",
			body: `{"team": {}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.TeamAdd(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestTeamAdd_TeamExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	teamInput := models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{User_id: "1", Username: "alice"},
		},
	}

	mockService.
		EXPECT().
		CreateOrUpdateTeam(&teamInput).
		Return(models.Team{}, errors.New("team_name already exists"))

	h := NewHandler(mockService, slog.Default())

	body, _ := json.Marshal(models.TeamAddRequest{Team: teamInput})
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TeamAdd(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTeamGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	testTeamName := "backend"

	expectedTeam := models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{
				User_id:  "u1",
				Username: "Alice",
				IsActive: true,
			},
			{
				User_id:  "u2",
				Username: "Bob",
				IsActive: true,
			},
		},
	}
	mockService.
		EXPECT().
		GetTeamWithMembers(testTeamName).
		Return(expectedTeam, nil)

	h := NewHandler(mockService, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+testTeamName, nil)
	w := httptest.NewRecorder()

	h.TeamGet(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out models.TeamAddResponse200
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	require.Equal(t, out.Team, expectedTeam)
}

func TestTeamGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTeamName := "unknown"

	mockSvc := NewMockService(ctrl)
	// Expect call with "unknown" and return error
	mockSvc.
		EXPECT().
		GetTeamWithMembers(testTeamName).
		Return(models.Team{}, errors.New("not found"))

	h := NewHandler(mockSvc, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+testTeamName, nil)
	w := httptest.NewRecorder()

	h.TeamGet(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestTeamGet_MissingTeamName(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	w := httptest.NewRecorder()

	h.TeamGet(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
func TestUsersSetIsActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	reqBody := models.UsersSetIsActiveRequest{
		UserId:   "u2",
		IsActive: false,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	expectedUser := models.User{
		UserId:   "u2",
		Username: "Bob",
		TeamName: "backend",
		IsActive: false,
	}

	mockService.EXPECT().UsersSetIsActive(reqBody.UserId, reqBody.IsActive).Return(expectedUser, nil)

	h.UsersSetIsActive(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp models.UsersSerIsActiveResponse200
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, expectedUser, resp.User)
}

func TestUsersSetIsActive_InvalidInput(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	tests := []struct {
		name string
		body string
	}{
		{name: "invalid JSON", body: "{invalid json"},
		{name: "empty user_id", body: `{"user_id": "" , "is_active": true}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()

			h.UsersSetIsActive(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUsersSetIsActive_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	reqBody := models.UsersSetIsActiveRequest{
		UserId:   "nonexistent",
		IsActive: true,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockService.EXPECT().
		UsersSetIsActive(reqBody.UserId, reqBody.IsActive).
		Return(models.User{}, errors.New("user not found"))

	h.UsersSetIsActive(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestPullRequestCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	reqBody := models.PullRequestCreateRequest{
		PullRequestId:   "pr-1001",
		PullRequestName: "Add search",
		AuthorId:        "u1",
	}
	body, _ := json.Marshal(reqBody)

	expectedPR := models.PullRequest{
		PullRequestId:     "pr-1001",
		PullRequestName:   "Add search",
		AuthorId:          "u1",
		Status:            models.StatusOpen,
		AssignedReviewers: []string{"u2", "u3"},
	}

	mockService.
		EXPECT().
		PullRequestCreate(reqBody.ToPullRequest()).
		Return(expectedPR, nil)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp models.PullRequestCreateResponse201
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, expectedPR, resp.Pr)
}

func TestPullRequestCreate_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestCreate_InvalidInput(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty pr id",
			body: `{"pull_request_id": "", "pull_request_name": "Add", "author_id": "u1"}`,
		},
		{
			name: "empty pr name",
			body: `{"pull_request_id": "pr-1", "pull_request_name": "", "author_id": "u1"}`,
		},
		{
			name: "empty author id",
			body: `{"pull_request_id": "pr-1", "pull_request_name": "Add", "author_id": ""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.PullRequestCreate(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestPullRequestCreate_PRExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestCreateRequest{
		PullRequestId:   "pr-1001",
		PullRequestName: "Feature",
		AuthorId:        "u1",
	}

	body, _ := json.Marshal(reqBody)

	mockSvc.EXPECT().
		PullRequestCreate(reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("PR id already exists"))

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestCreate_AuthorOrTeamNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestCreateRequest{
		PullRequestId:   "pr-2000",
		PullRequestName: "Fix",
		AuthorId:        "u999",
	}

	body, _ := json.Marshal(reqBody)

	mockSvc.EXPECT().
		PullRequestCreate(reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("author not found"))

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestCreate_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestCreateRequest{
		PullRequestId:   "pr-x",
		PullRequestName: "Refactor",
		AuthorId:        "u1",
	}
	body, _ := json.Marshal(reqBody)

	mockSvc.EXPECT().
		PullRequestCreate(reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("unknown error"))

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
