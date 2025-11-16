package models

import "fmt"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

func (s Status) Validate() error {
	switch s {
	case StatusOpen, StatusMerged:
		return nil
	default:
		return fmt.Errorf("bad status: %s", s)
	}
}
