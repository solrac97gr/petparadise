package models

type Pet struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Species     string   `json:"species" db:"species"`
	Breed       string   `json:"breed" db:"breed"`
	Age         int      `json:"age" db:"age"`
	Description string   `json:"description" db:"description"`
	Status      Status   `json:"status" db:"status"`
	Created     string   `json:"created" db:"created"`
	Updated     string   `json:"updated" db:"updated"`
	Images      []string `json:"images" db:"images"`
}

// NewPet creates a new Pet instance
func NewPet(id, name, species, breed string, age int, description string, status Status, images []string) (*Pet, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	if species == "" {
		return nil, ErrInvalidSpecies
	}

	if age < 0 {
		return nil, ErrInvalidAge
	}

	return &Pet{
		ID:          id,
		Name:        name,
		Species:     species,
		Breed:       breed,
		Age:         age,
		Description: description,
		Status:      status,
		Images:      images,
	}, nil
}
