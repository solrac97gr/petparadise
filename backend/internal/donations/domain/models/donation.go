package models

type Donation struct {
	ID        string  `json:"id" db:"id"`
	UserID    string  `json:"user_id" db:"user_id"`
	Amount    float64 `json:"amount" db:"amount"`
	Status    Status  `json:"status" db:"status"`
	Created   string  `json:"created" db:"created"`
	Updated   string  `json:"updated" db:"updated"`
	Comment   string  `json:"comment" db:"comment"`
	Anonymous bool    `json:"anonymous" db:"anonymous"`
}

// NewDonation creates a new Donation instance
func NewDonation(id, userID string, amount float64, status Status, comment string, anonymous bool) (*Donation, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return &Donation{
		ID:        id,
		UserID:    userID,
		Amount:    amount,
		Status:    status,
		Comment:   comment,
		Anonymous: anonymous,
	}, nil
}
