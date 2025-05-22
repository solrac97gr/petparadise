package models

import "errors"

type Status string

var (
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidAmount = errors.New("invalid amount")
)

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusRefunded  Status = "refunded"
)

var (
	validStatuses = map[Status]struct{}{
		StatusPending:   {},
		StatusCompleted: {},
		StatusFailed:    {},
		StatusRefunded:  {},
	}
)

// String converts the Status to a string
func (s Status) String() string {
	return string(s)
}

// IsEquals checks if the status is equal to another status
func (s Status) IsEquals(in Status) bool {
	return s.String() == in.String()
}

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	_, ok := validStatuses[s]
	return ok
}
