package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (h *Handler) PullRequestCreate(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.PullRequestCreateRequest
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
	pr, err := h.service.PullRequestCreate(req.ToPullRequest())
	if err != nil {
		// author/team not found
		// PR is already exists
		h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
		return
	}
	// create response
	resp := models.PullRequestCreateResponse201{
		Pr: pr,
	}
	// send response
	h.sendJSON(w, http.StatusCreated, resp)
}
