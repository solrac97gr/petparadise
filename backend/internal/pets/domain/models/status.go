package models

import "errors"

type Status string

var (
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidName    = errors.New("invalid name")
	ErrInvalidSpecies = errors.New("invalid species")
	ErrInvalidAge     = errors.New("invalid age")
)

const (
	StatusAvailable   Status = "available"
	StatusAdopted     Status = "adopted"
	StatusInProcess   Status = "in_process"
	StatusUnavailable Status = "unavailable"
	StatusQuarantined Status = "quarantined"
	StatusMedicalCare Status = "medical_care"
)

var (
	validStatuses = map[Status]struct{}{
		StatusAvailable:   {},
		StatusAdopted:     {},
		StatusInProcess:   {},
		StatusUnavailable: {},
		StatusQuarantined: {},
		StatusMedicalCare: {},
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
