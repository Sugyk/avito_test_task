package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (h *Handler) PullRequestCreate(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.PullRequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// validate request
	if err := req.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// business logic
	pr, err := h.service.PullRequestCreate(r.Context(), req.ToPullRequest())
	if err != nil {
		// author/team not found
		if errors.Is(err, models.ErrAuthorNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
			return
		}
		// PR is already exists
		if errors.Is(err, models.ErrPRAlreadyExists) {
			h.sendError(w, http.StatusConflict, models.PrExistsErrorCode, err)
			return
		}
		h.logger.Error("internal error", "error", err.Error())
		h.sendError(w, http.StatusInternalServerError, models.InternalErrorCode, models.ErrInternalError)
		return
	}
	// create response
	resp := models.PullRequestCreateResponse201{
		Pr: *pr,
	}
	// send response
	h.sendJSON(w, http.StatusCreated, resp)
}

func (h *Handler) PullRequestMerge(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.PullRequestMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// validate request
	if err := req.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// business logic
	pr, err := h.service.PullRequestMerge(r.Context(), req.ToPullRequest())
	if err != nil {
		// pr not found
		if errors.Is(err, models.ErrPRNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
			return
		}
		// PR is already exists
		h.logger.Error("internal error", "error", err.Error())
		h.sendError(w, http.StatusInternalServerError, models.InternalErrorCode, models.ErrInternalError)
		return
	}
	// create response
	resp := models.PullRequestMergeResponse200{
		Pr: *pr,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}

func (h *Handler) PullRequestReassign(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.PullRequestReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// validate request
	if err := req.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// business logic
	pr, replacedBy, err := h.service.PullRequestReassign(r.Context(), req.PullRequestId, req.OldReviewerId)
	if err != nil {
		// 404
		// pr not found
		if errors.Is(err, models.ErrPRNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
			return
		}
		// user not found
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
			return
		}
		// 409
		// reassigning merged pr
		if errors.Is(err, models.ErrReassigningMergedPR) {
			h.sendError(w, http.StatusConflict, models.PrMergedErrorCode, err)
			return
		}
		// not assigned user
		if errors.Is(err, models.ErrUserNotAssignedToPR) {
			h.sendError(w, http.StatusConflict, models.NotAssignedErrorCode, err)
			return
		}
		// no candidates
		if errors.Is(err, models.ErrNoActiveCandidates) {
			h.sendError(w, http.StatusConflict, models.NoCandidateErrorCode, err)
			return
		}
		h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
		return
	}
	// create response
	resp := models.PullRequestReassignResponse200{
		Pr:         *pr,
		ReplacedBy: replacedBy,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
