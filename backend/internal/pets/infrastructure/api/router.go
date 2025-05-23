package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/pets/aplication"
	"github.com/solrac97gr/petparadise/internal/pets/infrastructure/repository"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/pkg/auth"
)

// SetupPetRoutes sets up all pet routes
func SetupPetRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	petRepo := repository.NewPostgresRepository(db)

	// Initialize service
	petService := aplication.NewPetService(petRepo)

	// Initialize handler
	petHandler := NewPetHandler(petService)

	// Public routes - anyone can view pets
	router.Get("/", petHandler.GetAllPets)
	router.Get("/:id", petHandler.GetPetByID)
	router.Get("/status", petHandler.GetPetsByStatus)

	// Protected routes - require authentication
	protectedRoutes := router.Use(auth.Protected())

	// Staff routes - require authentication + proper role
	staffRoutes := protectedRoutes.Use(auth.RoleRequired(models.RoleAdmin, models.RoleVet, models.RoleVolunteer))
	staffRoutes.Post("/", petHandler.CreatePet)
	staffRoutes.Put("/:id", petHandler.UpdatePet)
	staffRoutes.Patch("/:id/status", petHandler.UpdatePetStatus)
	staffRoutes.Delete("/:id", petHandler.DeletePet)
}
