package ports

import "github.com/solrac97gr/petparadise/internal/pets/domain/models"

type PetRepository interface {
	Save(pet *models.Pet) error
	FindByID(id string) (*models.Pet, error)
	FindByStatus(status models.Status) ([]*models.Pet, error)
	FindAll() ([]*models.Pet, error)
	Update(pet *models.Pet) error
	Delete(id string) error
}

type PetService interface {
	CreatePet(name, species, breed string, age int, description string, images []string) (*models.Pet, error)
	GetPetByID(id string) (*models.Pet, error)
	GetPetsByStatus(status models.Status) ([]*models.Pet, error)
	GetAllPets() ([]*models.Pet, error)
	UpdatePet(id, name, species, breed string, age int, description string, status models.Status, images []string) (*models.Pet, error)
	UpdatePetStatus(id string, status models.Status) (*models.Pet, error)
	DeletePet(id string) error
}
