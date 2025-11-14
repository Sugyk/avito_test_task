package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/models"
)

// /users/setIsActive:
//
//	post:
//	  tags: [Users]
//	  summary: Установить флаг активности пользователя
//	  requestBody:
//	    required: true
//	    content:
//	      application/json:
//	        schema:
//	          type: object
//	          required: [ user_id, is_active ]
//	          properties:
//	            user_id:
//	              type: string
//	            is_active:
//	              type: boolean
//	        example:
//	          user_id: u2
//	          is_active: false
//	  responses:
//	    '200':
//	      description: Обновлённый пользователь
//	      content:
//	        application/json:
//	          schema:
//	            type: object
//	            properties:
//	              user:
//	                $ref: '#/components/schemas/User'
//	          example:
//	            user:
//	              user_id: u2
//	              username: Bob
//	              team_name: backend
//	              is_active: false
//	    '404':
//	      description: Пользователь не найден
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
func (h *Handler) UsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.UsersSetIsActiveRequest
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
	user, err := h.service.UsersSetIsActive(req.UserId, req.IsActive)
	if err != nil {
		// user not found
		h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
		return
	}
	// create response
	resp := models.UsersSerIsActiveResponse200{
		User: user,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
