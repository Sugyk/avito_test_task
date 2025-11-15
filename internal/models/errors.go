package models

import "errors"

var (
	TeamExistsErrorCode   = "TEAM_EXISTS"
	PrExistsErrorCode     = "PR_EXISTS"
	PrMergedErrorCode     = "PR_MERGED"
	NotAssignedErrorCode  = "NOT_ASSIGNED"
	NoCandidateErrorCode  = "NO_CANDIDATE"
	NotFoundErrorCode     = "NOT_FOUND"
	InvalidInputErrorCode = "INVALID_INPUT"
	InternalErrorCode     = "INTERNAL_ERROR"
)

var (
	ErrTeamExists    = errors.New("team_name already exists")
	ErrInternalError = errors.New("internal server error")
	ErrTeamNotFound  = errors.New("team_name not found")
	ErrUserNotFound  = errors.New("user not found")
)

type Error struct {
	Code    string
	Message string
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
