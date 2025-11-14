package models

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

func (s Status) Validate() bool {
	switch s {
	case StatusOpen, StatusMerged:
		return true
	default:
		return false
	}
}
