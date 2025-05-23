package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/adoptions/aplication"
	"github.com/solrac97gr/petparadise/internal/adoptions/infrastructure/repository"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/pkg/auth"
)

// SetupAdoptionRoutes sets up all adoption routes
func SetupAdoptionRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	adoptionRepo := repository.NewPostgresRepository(db)

	// Initialize service
	adoptionService := aplication.NewAdoptionService(adoptionRepo)

	// Initialize handler
	adoptionHandler := NewAdoptionHandler(adoptionService)

	// All adoption routes require authentication
	protected := router.Use(auth.Protected())

	// Regular user routes - users can create adoptions and see their own
	protected.Post("/", adoptionHandler.CreateAdoption)
	protected.Get("/:id", adoptionHandler.GetAdoptionByID)
	protected.Get("/user/:userId", adoptionHandler.GetAdoptionsByUserID)

	// Staff routes - require admin, volunteer or vet role
	staffRoutes := protected.Use(auth.RoleRequired(models.RoleAdmin, models.RoleVet, models.RoleVolunteer))
	staffRoutes.Get("/", adoptionHandler.GetAllAdoptions)
	staffRoutes.Put("/:id", adoptionHandler.UpdateAdoption)
	staffRoutes.Delete("/:id", adoptionHandler.DeleteAdoption)
}
