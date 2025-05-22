package ports

import "github.com/solrac97gr/petparadise/internal/adoptions/domain/models"

type AdoptionRepository interface {
	Save(adoption *models.Adoption) error
	FindByID(id string) (*models.Adoption, error)
	FindByUserID(userID string) ([]*models.Adoption, error)
	FindAll() ([]*models.Adoption, error)
	Update(adoption *models.Adoption) error
	Delete(id string) error
}

type AdoptionService interface {
	CreateAdoption(petID, userID string, documents []string) (*models.Adoption, error)
	GetAdoptionByID(id string) (*models.Adoption, error)
	GetAdoptionsByUserID(userID string) ([]*models.Adoption, error)
	GetAllAdoptions() ([]*models.Adoption, error)
	UpdateAdoption(id string, status models.Status, documents []string) (*models.Adoption, error)
	DeleteAdoption(id string) error
}
