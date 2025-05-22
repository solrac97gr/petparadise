package ports

import "github.com/solrac97gr/petparadise/internal/donations/domain/models"

type DonationRepository interface {
	Save(donation *models.Donation) error
	FindByID(id string) (*models.Donation, error)
	FindByUserID(userID string) ([]*models.Donation, error)
	FindAll() ([]*models.Donation, error)
	Update(donation *models.Donation) error
	Delete(id string) error
}

type DonationService interface {
	CreateDonation(userID string, amount float64, comment string, anonymous bool) (*models.Donation, error)
	GetDonationByID(id string) (*models.Donation, error)
	GetDonationsByUserID(userID string) ([]*models.Donation, error)
	GetAllDonations() ([]*models.Donation, error)
	UpdateDonation(id string, status models.Status) (*models.Donation, error)
	DeleteDonation(id string) error
}
