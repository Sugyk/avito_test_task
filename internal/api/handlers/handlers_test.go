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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func bool_pointer(x bool) *bool {
	return &x
}
func TestTeamAdd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	teamInput := &models.Team{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserId: "1", Username: "alice", IsActive: bool_pointer(true)},
		},
	}
	body, _ := json.Marshal(teamInput)
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

	var out models.TeamAddResponse201
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

	var expectedBody = models.ErrorResponse{
		Error: models.Error{
			Code:    models.InvalidInputErrorCode,
			Message: "invalid character 'i' looking for beginning of object key string",
		},
	}
	var outBody models.ErrorResponse
	err := json.NewDecoder(w.Result().Body).Decode(&outBody)
	require.NoError(t, err)
	require.Equal(t, expectedBody, outBody)

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
			{UserId: "1", Username: "alice", IsActive: bool_pointer(true)},
		},
	}
	body, _ := json.Marshal(teamInput)
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))

	mockService.
		EXPECT().
		CreateOrUpdateTeam(req.Context(), teamInput).
		Return(nil, models.ErrTeamExists)

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
				UserId:   "u1",
				Username: "Alice",
				IsActive: bool_pointer(true),
			},
			{
				UserId:   "u2",
				Username: "Bob",
				IsActive: bool_pointer(true),
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

	var out models.TeamAddResponse201
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
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestTeamGet_MissingTeamName(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	w := httptest.NewRecorder()

	h.TeamGet(w, req)

	resp := w.Result()
	defer func() {
		_ = resp.Body.Close()
	}()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
func TestUsersSetIsActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	h := NewHandler(mockService, slog.Default())

	reqBody := models.UsersSetIsActiveRequest{
		UserId:   "u2",
		IsActive: bool_pointer(false),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	expectedUser := &models.User{
		UserId:   "u2",
		Username: "Bob",
		TeamName: "backend",
		IsActive: false,
	}

	mockService.EXPECT().UsersSetIsActive(req.Context(), reqBody.UserId, *reqBody.IsActive).Return(expectedUser, nil)

	h.UsersSetIsActive(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp models.UsersSerIsActiveResponse200
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, *expectedUser, resp.User)
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
		IsActive: bool_pointer(true),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockService.EXPECT().
		UsersSetIsActive(req.Context(), reqBody.UserId, *reqBody.IsActive).
		Return(nil, models.ErrUserNotFound)

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

	expectedPR := &models.PullRequest{
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
	require.Equal(t, *expectedPR, resp.Pr)
}

func TestPullRequestCreate_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	var expectedBody = models.ErrorResponse{
		Error: models.Error{
			Code:    models.InvalidInputErrorCode,
			Message: "invalid character 'i' looking for beginning of object key string",
		},
	}
	var outBody models.ErrorResponse
	err := json.NewDecoder(w.Result().Body).Decode(&outBody)
	require.NoError(t, err)
	require.Equal(t, expectedBody, outBody)

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
		Return(nil, models.ErrPRAlreadyExists)

	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
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
		Return(nil, models.ErrAuthorNotFound)

	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestPullRequestCreate_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	reqBody := &models.PullRequestCreateRequest{
		PullRequestId:   "pr-x",
		PullRequestName: "Refactor",
		AuthorId:        "u1",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
	mockSvc.EXPECT().
		PullRequestCreate(req.Context(), reqBody.ToPullRequest()).
		Return(nil, errors.New("unknown error"))

	w := httptest.NewRecorder()

	h.PullRequestCreate(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
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

	expectedPR := &models.PullRequest{
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
	require.Equal(t, *expectedPR, resp.Pr)
}

func TestPullRequestMerge_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestMerge(w, req)

	var expectedBody = models.ErrorResponse{
		Error: models.Error{
			Code:    models.InvalidInputErrorCode,
			Message: "invalid character 'i' looking for beginning of object key string",
		},
	}
	var outBody models.ErrorResponse
	err := json.NewDecoder(w.Result().Body).Decode(&outBody)
	require.NoError(t, err)
	require.Equal(t, expectedBody, outBody)

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
		Return(nil, models.ErrPRNotFound)

	w := httptest.NewRecorder()

	h.PullRequestMerge(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestPullRequestReassign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockService(ctrl)
	h := NewHandler(mockSvc, slog.Default())

	testCases := []struct {
		name          string
		reqBody       models.PullRequestReassignRequest
		mockSetup     func()
		expectedCode  int
		expectedError *models.Error
	}{
		{
			name: "success",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-1",
				OldReviewerId: "user-1",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-1", "user-1").
					Return(
						&models.PullRequest{PullRequestId: "pr-1"},
						"user-2",
						nil,
					)
			},
			expectedCode:  http.StatusOK,
			expectedError: nil,
		},
		{
			name: "pull_request_id is missed",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "",
				OldReviewerId: "user-1",
			},
			mockSetup:    nil,
			expectedCode: http.StatusBadRequest,
			expectedError: &models.Error{
				Code:    models.InvalidInputErrorCode,
				Message: "pull_request_id is required",
			},
		},
		{
			name: "old_reviewer_id is missed",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-1",
				OldReviewerId: "",
			},
			mockSetup:    nil,
			expectedCode: http.StatusBadRequest,
			expectedError: &models.Error{
				Code:    models.InvalidInputErrorCode,
				Message: "old_reviewer_id is required",
			},
		},
		{
			name: "PR not found",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-999",
				OldReviewerId: "user-1",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-999", "user-1").
					Return(nil, "", models.ErrPRNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedError: &models.Error{
				Code:    models.NotFoundErrorCode,
				Message: models.ErrPRNotFound.Error(),
			},
		},
		{
			name: "User not found",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-1",
				OldReviewerId: "user-999",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-1", "user-999").
					Return(nil, "", models.ErrUserNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedError: &models.Error{
				Code:    models.NotFoundErrorCode,
				Message: models.ErrUserNotFound.Error(),
			},
		},
		{
			name: "Попытка переназначения в merged PR",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-2",
				OldReviewerId: "user-1",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-2", "user-1").
					Return(nil, "", models.ErrReassigningMergedPR)
			},
			expectedCode: http.StatusConflict,
			expectedError: &models.Error{
				Code:    models.PrMergedErrorCode,
				Message: models.ErrReassigningMergedPR.Error(),
			},
		},
		{
			name: "Пользователь не назначен на PR",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-3",
				OldReviewerId: "user-3",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-3", "user-3").
					Return(nil, "", models.ErrUserNotAssignedToPR)
			},
			expectedCode: http.StatusConflict,
			expectedError: &models.Error{
				Code:    models.NotAssignedErrorCode,
				Message: models.ErrUserNotAssignedToPR.Error(),
			},
		},
		{
			name: "Нет активных кандидатов",
			reqBody: models.PullRequestReassignRequest{
				PullRequestId: "pr-4",
				OldReviewerId: "user-1",
			},
			mockSetup: func() {
				mockSvc.EXPECT().
					PullRequestReassign(gomock.Any(), "pr-4", "user-1").
					Return(nil, "", models.ErrNoActiveCandidates)
			},
			expectedCode: http.StatusConflict,
			expectedError: &models.Error{
				Code:    models.NoCandidateErrorCode,
				Message: models.ErrNoActiveCandidates.Error(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockSetup != nil {
				tc.mockSetup()
			}

			body, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest(
				http.MethodPost,
				"/pull-request/reassign",
				bytes.NewReader(body),
			)

			w := httptest.NewRecorder()

			h.PullRequestReassign(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != nil {
				var resp models.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedError.Code, resp.Error.Code)
				assert.Equal(t, tc.expectedError.Message, resp.Error.Message)
			} else {
				var resp interface{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
			}
		})
	}
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

	expectedPR := &models.PullRequest{
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
	require.Equal(t, *expectedPR, resp.Pr)
	require.Equal(t, "u5", resp.ReplacedBy)
}

func TestPullRequestReassign_InvalidJSON(t *testing.T) {
	h := NewHandler(nil, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	h.PullRequestReassign(w, req)

	var expectedBody = models.ErrorResponse{
		Error: models.Error{
			Code:    models.InvalidInputErrorCode,
			Message: "invalid character 'i' looking for beginning of object key string",
		},
	}
	var outBody models.ErrorResponse
	err := json.NewDecoder(w.Result().Body).Decode(&outBody)
	require.NoError(t, err)
	require.Equal(t, expectedBody, outBody)

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
	cases := []struct {
		name   string
		svcErr error
	}{
		{"pr not found", models.ErrPRNotFound},
		{"no candidate", errors.New("no active replacement candidate in team")},
	}

	reqBody := &models.PullRequestReassignRequest{
		PullRequestId: "unknown-pr",
		OldReviewerId: "u777",
	}

	body, _ := json.Marshal(reqBody)
	for _, c := range cases {
		t.Run(
			c.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
				mockSvc.
					EXPECT().
					PullRequestReassign(req.Context(), reqBody.PullRequestId, reqBody.OldReviewerId).
					Return(nil, "", models.ErrPRNotFound)

				w := httptest.NewRecorder()

				h.PullRequestReassign(w, req)

				require.Equal(t, http.StatusNotFound, w.Code)
			},
		)
	}
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
		{"PR merged", models.ErrReassigningMergedPR},
		{"not assigned", models.ErrUserNotAssignedToPR},
		{"no candidate", models.ErrNoActiveCandidates},
	}

	reqBody := models.PullRequestReassignRequest{
		PullRequestId: "pr-x",
		OldReviewerId: "u5",
	}
	body, _ := json.Marshal(reqBody)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))

			mockSvc.
				EXPECT().
				PullRequestReassign(req.Context(), reqBody.PullRequestId, reqBody.OldReviewerId).
				Return(nil, "", c.svcErr)

			w := httptest.NewRecorder()

			h.PullRequestReassign(w, req)

			require.Equal(t, http.StatusConflict, w.Code)
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
