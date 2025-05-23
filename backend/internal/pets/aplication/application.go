package aplication

import (
	"time"

	"github.com/google/uuid"
	"github.com/solrac97gr/petparadise/internal/pets/domain/models"
	"github.com/solrac97gr/petparadise/internal/pets/domain/ports"
)

type PetService struct {
	repository ports.PetRepository
}

// NewPetService creates a new PetService instance
func NewPetService(repository ports.PetRepository) *PetService {
	return &PetService{
		repository: repository,
	}
}

// CreatePet creates a new pet
func (s *PetService) CreatePet(name, species, breed string, age int, description string, images []string) (*models.Pet, error) {
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	pet, err := models.NewPet(id, name, species, breed, age, description, models.StatusAvailable, images)
	if err != nil {
		return nil, err
	}

	pet.Created = now
	pet.Updated = now

	err = s.repository.Save(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

// GetPetByID returns a pet by its ID
func (s *PetService) GetPetByID(id string) (*models.Pet, error) {
	return s.repository.FindByID(id)
}

// GetPetsByStatus returns all pets with a specific status
func (s *PetService) GetPetsByStatus(status models.Status) ([]*models.Pet, error) {
	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}
	return s.repository.FindByStatus(status)
}

// GetAllPets returns all pets
func (s *PetService) GetAllPets() ([]*models.Pet, error) {
	return s.repository.FindAll()
}

// UpdatePet updates a pet's information
func (s *PetService) UpdatePet(id, name, species, breed string, age int, description string, status models.Status, images []string) (*models.Pet, error) {
	pet, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if pet == nil {
		return nil, models.ErrInvalidStatus // In a real app, you'd have a more specific error like ErrPetNotFound
	}

	if name != "" {
		pet.Name = name
	}

	if species != "" {
		pet.Species = species
	}

	pet.Breed = breed

	if age >= 0 {
		pet.Age = age
	}

	pet.Description = description

	if status != "" && !status.IsValid() {
		return nil, models.ErrInvalidStatus
	} else if status != "" {
		pet.Status = status
	}

	if images != nil {
		pet.Images = images
	}

	pet.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

// UpdatePetStatus updates only a pet's status
func (s *PetService) UpdatePetStatus(id string, status models.Status) (*models.Pet, error) {
	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}

	pet, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if pet == nil {
		return nil, models.ErrInvalidStatus // In a real app, you'd have a more specific error like ErrPetNotFound
	}

	pet.Status = status
	pet.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

// DeletePet deletes a pet
func (s *PetService) DeletePet(id string) error {
	return s.repository.Delete(id)
}
