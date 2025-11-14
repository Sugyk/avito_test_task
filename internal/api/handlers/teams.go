package handlers

import (
	"encoding/json"
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
	team, err := h.service.CreateOrUpdateTeam(&req.Team)
	if err != nil {
		// team_name already exists
		h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
		return
	}
	// create response
	resp := models.TeamAddResponse200{
		Team: team,
	}
	// send response
	h.sendJSON(w, http.StatusCreated, resp)
}
