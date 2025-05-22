package models

type Adoption struct {
	ID        string   `json:"id" db:"id"`
	PetID     string   `json:"pet_id" db:"pet_id"`
	UserID    string   `json:"user_id" db:"user_id"`
	Status    Status   `json:"status" db:"status"`
	Created   string   `json:"created" db:"created"`
	Updated   string   `json:"updated" db:"updated"`
	Documents []string `json:"documents" db:"documents"`
}

// NewAdoption creates a new Adoption instance
func NewAdoption(id, petID, userID string, status Status, documents []string) (*Adoption, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	return &Adoption{
		ID:        id,
		PetID:     petID,
		UserID:    userID,
		Status:    status,
		Documents: documents,
	}, nil
}
