package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (h *Handler) UsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.UsersSetIsActiveRequest
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
	user, err := h.service.UsersSetIsActive(r.Context(), req.UserId, *req.IsActive)
	if err != nil {
		// user not found
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
			return
		}
		h.logger.Error("internal error", "error", err.Error())
		h.sendError(w, http.StatusInternalServerError, models.InternalErrorCode, models.ErrInternalError)
		return
	}
	// create response
	resp := models.UsersSerIsActiveResponse200{
		User: *user,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}

func (h *Handler) UsersGetReview(w http.ResponseWriter, r *http.Request) {
	// extract query params
	userID := r.URL.Query().Get("user_id")
	// validate params
	if userID == "" {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, errors.New("missing user_id"))
		return
	}
	// business logic
	prs, err := h.service.UsersGetReview(r.Context(), userID)
	if err != nil {
		// user not found
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, models.ErrUserNotFound)
			return
		}
		h.logger.Error("error getting user's reviews", "error", err.Error())
		h.sendError(w, http.StatusInternalServerError, models.InternalErrorCode, models.ErrInternalError)
	}
	// create response
	resp := models.UsersGetReviewResponse200{
		UserId:       userID,
		PullRequests: prs,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
