package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (h *Handler) TeamAdd(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	// validate request
	if err := req.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// business logic
	team, err := h.service.CreateOrUpdateTeam(r.Context(), &req.Team)
	if err != nil {
		// team_name already exists
		if errors.Is(err, models.ErrTeamExists) {
			h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
			return
		}
		if errors.Is(err, models.ErrInternalError) {
			h.sendError(w, http.StatusInternalServerError, models.TeamExistsErrorCode, err)
			return
		}
	}
	// create response
	resp := models.TeamAddResponse200{
		Team: *team,
	}
	// send response
	h.sendJSON(w, http.StatusCreated, resp)
}

func (h *Handler) TeamGet(w http.ResponseWriter, r *http.Request) {
	// extract query params
	teamName := r.URL.Query().Get("team_name")
	// validate params
	if teamName == "" {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, errors.New("missing team_name"))
		return
	}
	// business logic
	team, err := h.service.GetTeamWithMembers(r.Context(), teamName)
	if err != nil {
		// team not found
		h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
		return
	}
	// create response
	resp := models.TeamGetResponse200{
		Team: team,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
