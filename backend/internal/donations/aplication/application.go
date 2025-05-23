package aplication

import (
	"time"

	"github.com/google/uuid"
	"github.com/solrac97gr/petparadise/internal/donations/domain/models"
	"github.com/solrac97gr/petparadise/internal/donations/domain/ports"
)

type DonationService struct {
	repository ports.DonationRepository
}

// NewDonationService creates a new DonationService instance
func NewDonationService(repository ports.DonationRepository) *DonationService {
	return &DonationService{
		repository: repository,
	}
}

// CreateDonation creates a new donation
func (s *DonationService) CreateDonation(userID string, amount float64, comment string, anonymous bool) (*models.Donation, error) {
	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	donation, err := models.NewDonation(id, userID, amount, models.StatusPending, comment, anonymous)
	if err != nil {
		return nil, err
	}

	donation.Created = now
	donation.Updated = now

	err = s.repository.Save(donation)
	if err != nil {
		return nil, err
	}

	return donation, nil
}

// GetDonationByID returns a donation by its ID
func (s *DonationService) GetDonationByID(id string) (*models.Donation, error) {
	return s.repository.FindByID(id)
}

// GetDonationsByUserID returns all donations for a user
func (s *DonationService) GetDonationsByUserID(userID string) ([]*models.Donation, error) {
	return s.repository.FindByUserID(userID)
}

// GetAllDonations returns all donations
func (s *DonationService) GetAllDonations() ([]*models.Donation, error) {
	return s.repository.FindAll()
}

// UpdateDonation updates a donation's status
func (s *DonationService) UpdateDonation(id string, status models.Status) (*models.Donation, error) {
	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}

	donation, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if donation == nil {
		return nil, models.ErrInvalidStatus // In a real app, you'd have a more specific error like ErrDonationNotFound
	}

	donation.Status = status
	donation.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(donation)
	if err != nil {
		return nil, err
	}

	return donation, nil
}

// DeleteDonation deletes a donation
func (s *DonationService) DeleteDonation(id string) error {
	return s.repository.Delete(id)
}
