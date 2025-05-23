package aplication

import (
	"time"

	"github.com/google/uuid"
	"github.com/solrac97gr/petparadise/internal/adoptions/domain/models"
	"github.com/solrac97gr/petparadise/internal/adoptions/domain/ports"
)

type AdoptionService struct {
	repository ports.AdoptionRepository
}

// NewAdoptionService creates a new AdoptionService instance
func NewAdoptionService(repository ports.AdoptionRepository) *AdoptionService {
	return &AdoptionService{
		repository: repository,
	}
}

// CreateAdoption creates a new adoption request
func (s *AdoptionService) CreateAdoption(petID, userID string, documents []string) (*models.Adoption, error) {
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	adoption, err := models.NewAdoption(id, petID, userID, models.StatusPending, documents)
	if err != nil {
		return nil, err
	}

	adoption.Created = now
	adoption.Updated = now

	err = s.repository.Save(adoption)
	if err != nil {
		return nil, err
	}

	return adoption, nil
}

// GetAdoptionByID returns an adoption by its ID
func (s *AdoptionService) GetAdoptionByID(id string) (*models.Adoption, error) {
	return s.repository.FindByID(id)
}

// GetAdoptionsByUserID returns all adoptions for a user
func (s *AdoptionService) GetAdoptionsByUserID(userID string) ([]*models.Adoption, error) {
	return s.repository.FindByUserID(userID)
}

// GetAllAdoptions returns all adoptions
func (s *AdoptionService) GetAllAdoptions() ([]*models.Adoption, error) {
	return s.repository.FindAll()
}

// UpdateAdoption updates an adoption
func (s *AdoptionService) UpdateAdoption(id string, status models.Status, documents []string) (*models.Adoption, error) {
	adoption, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}

	adoption.Status = status

	if documents != nil {
		adoption.Documents = documents
	}

	adoption.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(adoption)
	if err != nil {
		return nil, err
	}

	return adoption, nil
}

// DeleteAdoption deletes an adoption
func (s *AdoptionService) DeleteAdoption(id string) error {
	return s.repository.Delete(id)
}
