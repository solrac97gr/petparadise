package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/donations/aplication"
	"github.com/solrac97gr/petparadise/internal/donations/infrastructure/repository"
)

// SetupDonationRoutes sets up all donation routes
func SetupDonationRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	donationRepo := repository.NewPostgresRepository(db)

	// Initialize service
	donationService := aplication.NewDonationService(donationRepo)

	// Initialize handler
	donationHandler := NewDonationHandler(donationService)

	// Donation routes
	router.Post("/", donationHandler.CreateDonation)
	router.Get("/", donationHandler.GetAllDonations)
	router.Get("/:id", donationHandler.GetDonationByID)
	router.Get("/user/:userId", donationHandler.GetDonationsByUserID)
	router.Patch("/:id/status", donationHandler.UpdateDonationStatus)
	router.Delete("/:id", donationHandler.DeleteDonation)
}
