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
		h.sendError(w, http.StatusBadRequest, "invalid request body", err)
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
		h.sendError(w, http.StatusBadRequest, "invalid request body", err)
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
		// PR is already exists
		h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
		return
	}
	// create response
	resp := models.PullRequestMergeResponse200{
		Pr: pr,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}

// /pullRequest/reassign:
//
//	post:
//	  tags: [PullRequests]
//	  summary: Переназначить конкретного ревьювера на другого из его команды
//	  requestBody:
//	    required: true
//	    content:
//	      application/json:
//	        schema:
//	          type: object
//	          required: [ pull_request_id, old_user_id ]
//	          properties:
//	            pull_request_id: { type: string }
//	            old_user_id: { type: string }
//	        example:
//	          pull_request_id: pr-1001
//	          old_reviewer_id: u2
//	  responses:
//	    '200':
//	      description: Переназначение выполнено
//	      content:
//	        application/json:
//	          schema:
//	            type: object
//	            required: [pr, replaced_by]
//	            properties:
//	              pr:
//	                $ref: '#/components/schemas/PullRequest'
//	              replaced_by:
//	                type: string
//	                description: user_id нового ревьювера
//	          example:
//	            pr:
//	              pull_request_id: pr-1001
//	              pull_request_name: Add search
//	              author_id: u1
//	              status: OPEN
//	              assigned_reviewers: [u3, u5]
//	            replaced_by: u5
//	    '404':
//	      description: PR или пользователь не найден
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
//	    '409':
//	      description: Нарушение доменных правил переназначения
//	      content:
//	        application/json:
//	          schema: { $ref: '#/components/schemas/ErrorResponse' }
//	          examples:
//	            merged:
//	              summary: Нельзя менять после MERGED
//	              value:
//	                error: { code: PR_MERGED, message: cannot reassign on merged PR }
//	            notAssigned:
//	              summary: Пользователь не был назначен ревьювером
//	              value:
//	                error: { code: NOT_ASSIGNED, message: reviewer is not assigned to this PR }
//	            noCandidate:
//	              summary: Нет доступных кандидатов
//	              value:
//	                error: { code: NO_CANDIDATE, message: no active replacement candidate in team }
func (h *Handler) PullRequestReassign(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req models.PullRequestReassignRequest
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
	pr, replacedBy, err := h.service.PullRequestReassign(r.Context(), req.PullRequestId, req.OldReviewerId)
	if err != nil {
		// pr not found
		// PR is already exists
		h.sendError(w, http.StatusBadRequest, models.TeamExistsErrorCode, err)
		return
	}
	// create response
	resp := models.PullRequestReassignResponse200{
		Pr:         pr,
		ReplacedBy: replacedBy,
	}
	// send response
	h.sendJSON(w, http.StatusOK, resp)
}
