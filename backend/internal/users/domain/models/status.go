package models

import "errors"

type Status string

var (
	ErrInvalidStatus   = errors.New("invalid status")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
	StatusPending   Status = "pending"
)

var (
	validStatuses = map[Status]struct{}{
		StatusActive:    {},
		StatusInactive:  {},
		StatusSuspended: {},
		StatusPending:   {},
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
