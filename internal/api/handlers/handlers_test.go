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

	teamInput := &models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{User_id: "1", Username: "alice"},
		},
	}
	body, _ := json.Marshal(models.TeamAddRequest{Team: *teamInput})
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))

	mockService.
		EXPECT().
		CreateOrUpdateTeam(req.Context(), teamInput).
		Return(teamInput, nil)

	h := NewHandler(mockService, nil)

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

	teamInput := &models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{User_id: "1", Username: "alice"},
		},
	}
	body, _ := json.Marshal(models.TeamAddRequest{Team: *teamInput})
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))

	mockService.
		EXPECT().
		CreateOrUpdateTeam(req.Context(), teamInput).
		Return(&models.Team{}, errors.New("team_name already exists"))

	h := NewHandler(mockService, slog.Default())

	w := httptest.NewRecorder()

	h.TeamAdd(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTeamGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	testTeamName := "backend"

	expectedTeam := &models.Team{
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
	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+testTeamName, nil)
	mockService.
		EXPECT().
		GetTeamWithMembers(req.Context(), testTeamName).
		Return(expectedTeam, nil)

	h := NewHandler(mockService, slog.Default())

	w := httptest.NewRecorder()

	h.TeamGet(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out models.TeamAddResponse200
	err := json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	require.Equal(t, out.Team, *expectedTeam)
}

func TestTeamGet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTeamName := "unknown"

	mockSvc := NewMockService(ctrl)
	// Expect call with "unknown" and return error
	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+testTeamName, nil)
	mockSvc.
		EXPECT().
		GetTeamWithMembers(req.Context(), testTeamName).
		Return(&models.Team{}, errors.New("not found"))

	h := NewHandler(mockSvc, slog.Default())

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

	mockService.EXPECT().UsersSetIsActive(req.Context(), reqBody.UserId, reqBody.IsActive).Return(expectedUser, nil)

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
		UsersSetIsActive(req.Context(), reqBody.UserId, reqBody.IsActive).
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
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))

	mockService.
		EXPECT().
		PullRequestCreate(req.Context(), reqBody.ToPullRequest()).
		Return(expectedPR, nil)

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

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	mockSvc.EXPECT().
		PullRequestCreate(req.Context(), reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("PR id already exists"))

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

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	mockSvc.EXPECT().
		PullRequestCreate(req.Context(), reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("author not found"))

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

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	mockSvc.EXPECT().
		PullRequestCreate(req.Context(), reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("unknown error"))

	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestMerge_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	reqBody := models.PullRequestMergeRequest{
		PullRequestId: "pr-1001",
	}
	body, _ := json.Marshal(reqBody)

	mergedtime := func() *string {
		v := "2025-11-14T10:00:00Z"
		return &v
	}()

	expectedPR := models.PullRequest{
		PullRequestId:     "pr-1001",
		PullRequestName:   "Add search",
		AuthorId:          "u1",
		Status:            models.StatusMerged,
		AssignedReviewers: []string{"u2", "u3"},
		MergedAt:          mergedtime,
	}

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
	mockService.
		EXPECT().
		PullRequestMerge(req.Context(), reqBody.ToPullRequest()).
		Return(expectedPR, nil)

	w := httptest.NewRecorder()

	h.PullRequestMerge(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp models.PullRequestMergeResponse200
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, expectedPR, resp.Pr)
}

func TestPullRequestMerge_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestMerge(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestMerge_InvalidInput(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty pr id",
			body: `{"pull_request_id": ""}`,
		},
		{
			name: "not provided",
			body: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()

			h.PullRequestMerge(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestPullRequestMerge_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestMergeRequest{
		PullRequestId: "unknown-pr",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
	mockSvc.EXPECT().
		PullRequestMerge(req.Context(), reqBody.ToPullRequest()).
		Return(models.PullRequest{}, errors.New("PR not found"))

	w := httptest.NewRecorder()

	h.PullRequestMerge(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestReassign_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestReassignRequest{
		PullRequestId: "pr-1001",
		OldReviewerId: "u2",
	}

	body, _ := json.Marshal(reqBody)

	expectedPR := models.PullRequest{
		PullRequestId:     "pr-1001",
		PullRequestName:   "Add search",
		AuthorId:          "u1",
		Status:            models.StatusOpen,
		AssignedReviewers: []string{"u3", "u5"},
	}

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	mockSvc.
		EXPECT().
		PullRequestReassign(req.Context(), reqBody.PullRequestId, reqBody.OldReviewerId).
		Return(expectedPR, "u5", nil)

	w := httptest.NewRecorder()

	h.PullRequestReassign(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp models.PullRequestReassignResponse200
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, expectedPR, resp.Pr)
	require.Equal(t, "u5", resp.ReplacedBy)
}

func TestPullRequestReassign_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestReassign(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestReassign_InvalidInput(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing pull_request_id",
			body: `{"old_reviewer_id":"u1"}`,
		},
		{
			name: "missing old_reviewer_id",
			body: `{"pull_request_id":"pr-1"}`,
		},
		{
			name: "empty fields",
			body: `{"pull_request_id":"", "old_reviewer_id":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.PullRequestReassign(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestPullRequestReassign_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := models.PullRequestReassignRequest{
		PullRequestId: "unknown-pr",
		OldReviewerId: "u777",
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
	mockSvc.
		EXPECT().
		PullRequestReassign(req.Context(), reqBody.PullRequestId, reqBody.OldReviewerId).
		Return(models.PullRequest{}, "", errors.New("not found"))

	w := httptest.NewRecorder()

	h.PullRequestReassign(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPullRequestReassign_DomainErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	cases := []struct {
		name   string
		svcErr error
	}{
		{"PR merged", errors.New("cannot reassign on merged PR")},
		{"not assigned", errors.New("reviewer is not assigned to this PR")},
		{"no candidate", errors.New("no active replacement candidate in team")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			reqBody := models.PullRequestReassignRequest{
				PullRequestId: "pr-x",
				OldReviewerId: "u5",
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))

			mockSvc.
				EXPECT().
				PullRequestReassign(req.Context(), reqBody.PullRequestId, reqBody.OldReviewerId).
				Return(models.PullRequest{}, "", c.svcErr)

			w := httptest.NewRecorder()

			h.PullRequestReassign(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUsersGetReview_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	userID := "u2"
	prs := []models.PullRequestShort{
		{
			PullRequestId:   "pr-1001",
			PullRequestName: "Add search",
			AuthorId:        "u1",
			Status:          "OPEN",
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id="+userID, nil)

	// mock service call
	mockService.EXPECT().
		UsersGetReview(req.Context(), userID).
		Return(prs, nil)

	w := httptest.NewRecorder()

	h.UsersGetReview(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	expectedJSON := `{"user_id":"u2","pull_requests":[{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"u1","status":"OPEN"}]}`
	require.JSONEq(t, expectedJSON, w.Body.String())
}

func TestUsersGetReview_MissingUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)
	w := httptest.NewRecorder()

	h.UsersGetReview(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "missing user_id")
}
