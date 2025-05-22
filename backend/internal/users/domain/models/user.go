package models

type User struct {
	ID        string   `json:"id" db:"id"`
	Name      string   `json:"name" db:"name"`
	Email     string   `json:"email" db:"email"`
	Password  string   `json:"-" db:"password"`
	Status    Status   `json:"status" db:"status"`
	Created   string   `json:"created" db:"created"`
	Updated   string   `json:"updated" db:"updated"`
	Role      Role     `json:"role" db:"role"`
	Address   string   `json:"address" db:"address"`
	Phone     string   `json:"phone" db:"phone"`
	Documents []string `json:"documents" db:"documents"`
}

// NewUser creates a new User instance
func NewUser(id, name, email, password string, status Status, role Role, address, phone string, documents []string) (*User, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	if email == "" {
		return nil, ErrInvalidEmail
	}

	if password == "" {
		return nil, ErrInvalidPassword
	}

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		Status:    status,
		Role:      role,
		Address:   address,
		Phone:     phone,
		Documents: documents,
	}, nil
}
