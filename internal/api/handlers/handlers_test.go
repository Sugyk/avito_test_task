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
	req := httptest.NewRequest(http.MethodPost, "/teams/add", bytes.NewReader(body))
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

	req := httptest.NewRequest(http.MethodPost, "/teams/add", bytes.NewBufferString("{invalid json"))
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
			req := httptest.NewRequest(http.MethodPost, "/teams/add", bytes.NewBufferString(tt.body))
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
		Return(models.Team{}, errors.New("team already exists"))

	h := NewHandler(mockService, slog.Default())

	body, _ := json.Marshal(models.TeamAddRequest{Team: teamInput})
	req := httptest.NewRequest(http.MethodPost, "/teams/add", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.TeamAdd(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
