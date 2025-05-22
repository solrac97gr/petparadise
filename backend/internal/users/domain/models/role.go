package models

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleVolunteer Role = "volunteer"
	RoleVet       Role = "vet"
)

var (
	validRoles = map[Role]struct{}{
		RoleAdmin:     {},
		RoleUser:      {},
		RoleVolunteer: {},
		RoleVet:       {},
	}
)

// String converts the Role to a string
func (r Role) String() string {
	return string(r)
}

// IsEquals checks if the role is equal to another role
func (r Role) IsEquals(in Role) bool {
	return r.String() == in.String()
}

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	_, ok := validRoles[r]
	return ok
}
