package models

var (
	TeamExistsErrorCode   = "TEAM_EXISTS"
	PrExistsErrorCode     = "PR_EXISTS"
	PrMergedErrorCode     = "PR_MERGED"
	NotAssignedErrorCode  = "NOT_ASSIGNED"
	NoCandidateErrorCode  = "NO_CANDIDATE"
	NotFoundErrorCode     = "NOT_FOUND"
	InvalidInputErrorCode = "INVALID_INPUT"
)

type Error struct {
	Code    string
	Message string
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
