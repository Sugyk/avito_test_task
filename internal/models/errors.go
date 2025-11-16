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
	ErrTeamExists          = errors.New("team_name already exists")
	ErrInternalError       = errors.New("internal server error")
	ErrTeamNotFound        = errors.New("team_name not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrPRAlreadyExists     = errors.New("PR id already exists")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrPRNotFound          = errors.New("PR not found")
	ErrReassigningMergedPR = errors.New("cannot reassign on merged PR")
	ErrUserNotAssignedToPR = errors.New("reviewer is not assigned to this PR")
	ErrNoActiveCandidates  = errors.New("no active replacement candidate in team")
	ErrNoReviewers         = errors.New("no reviewers assigned to PR")
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
