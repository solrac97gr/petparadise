package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/donations/aplication"
	"github.com/solrac97gr/petparadise/internal/donations/infrastructure/repository"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/pkg/auth"
)

// SetupDonationRoutes sets up all donation routes
func SetupDonationRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	donationRepo := repository.NewPostgresRepository(db)

	// Initialize service
	donationService := aplication.NewDonationService(donationRepo)

	// Initialize handler
	donationHandler := NewDonationHandler(donationService)

	// All donation routes require authentication
	protected := router.Use(auth.Protected())

	// User routes - authenticated users can make donations and see their own
	protected.Post("/", donationHandler.CreateDonation)
	protected.Get("/user/:userId", donationHandler.GetDonationsByUserID)
	protected.Get("/:id", donationHandler.GetDonationByID)

	// Admin routes - only administrators can see all donations and modify them
	adminRoutes := protected.Use(auth.RoleRequired(models.RoleAdmin))
	adminRoutes.Get("/", donationHandler.GetAllDonations)
	adminRoutes.Patch("/:id/status", donationHandler.UpdateDonationStatus)
	adminRoutes.Delete("/:id", donationHandler.DeleteDonation)
}
