package models

import "errors"

type Status string

var (
	ErrInvalidStatus = errors.New("invalid status")
)

const (
	StatusPending             Status = "pending"
	StatusApproved            Status = "approved"
	StatusRejected            Status = "rejected"
	StatusCompleted           Status = "completed"
	StatusCancelled           Status = "cancelled"
	StatusWaitingForDocuments Status = "waiting_for_documents"
	StatusInProgress          Status = "in_progress"
)

var (
	validStatuses = map[Status]struct{}{
		StatusPending:             {},
		StatusApproved:            {},
		StatusRejected:            {},
		StatusCompleted:           {},
		StatusCancelled:           {},
		StatusInProgress:          {},
		StatusWaitingForDocuments: {},
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
