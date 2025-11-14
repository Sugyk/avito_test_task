package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

func (h *Handler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", "error", err.Error())
	}
}

func (h *Handler) sendError(w http.ResponseWriter, status int, code string, err error) {
	h.logger.Error("request error", code, err.Error())
	resp := models.ErrorResponse{
		Error: models.Error{
			Code:    code,
			Message: err.Error(),
		},
	}

	h.sendJSON(w, status, resp)
}
