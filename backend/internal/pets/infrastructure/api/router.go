package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/pets/aplication"
	"github.com/solrac97gr/petparadise/internal/pets/infrastructure/repository"
)

// SetupPetRoutes sets up all pet routes
func SetupPetRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	petRepo := repository.NewPostgresRepository(db)

	// Initialize service
	petService := aplication.NewPetService(petRepo)

	// Initialize handler
	petHandler := NewPetHandler(petService)

	// Pet routes
	router.Post("/", petHandler.CreatePet)
	router.Get("/", petHandler.GetAllPets)
	router.Get("/:id", petHandler.GetPetByID)
	router.Get("/status", petHandler.GetPetsByStatus)
	router.Put("/:id", petHandler.UpdatePet)
	router.Patch("/:id/status", petHandler.UpdatePetStatus)
	router.Delete("/:id", petHandler.DeletePet)
}
