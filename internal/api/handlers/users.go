package handlers

import (
	"encoding/json"
	"errors"
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
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// validate request
	if err := req.Validate(); err != nil {
		h.sendError(w, http.StatusBadRequest, models.InvalidInputErrorCode, err)
		return
	}
	// business logic
	user, err := h.service.UsersSetIsActive(r.Context(), req.UserId, req.IsActive)
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

//   /users/getReview:
//     get:
//       tags: [Users]
//       summary: Получить PR'ы, где пользователь назначен ревьювером
//       parameters:
//         - $ref: '#/components/parameters/UserIdQuery'
//       responses:
//         '200':
//           description: Список PR'ов пользователя
//           content:
//             application/json:
//               schema:
//                 type: object
//                 required: [ user_id, pull_requests ]
//                 properties:
//                   user_id:
//                     type: string
//                   pull_requests:
//                     type: array
//                     items:
//                       $ref: '#/components/schemas/PullRequestShort'
//               example:
//                 user_id: u2
//                 pull_requests:
//                   - pull_request_id: pr-1001
//                     pull_request_name: Add search
//                     author_id: u1
//                     status: OPEN

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
		h.sendError(w, http.StatusNotFound, models.NotFoundErrorCode, err)
		return
	}
	// create response
	resp := models.UsersGetReviewResponse200{
		UserId:       userID,
		PullRequests: prs,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
