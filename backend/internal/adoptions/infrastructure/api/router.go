package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/adoptions/aplication"
	"github.com/solrac97gr/petparadise/internal/adoptions/infrastructure/repository"
)

// SetupAdoptionRoutes sets up all adoption routes
func SetupAdoptionRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	adoptionRepo := repository.NewPostgresRepository(db)

	// Initialize service
	adoptionService := aplication.NewAdoptionService(adoptionRepo)

	// Initialize handler
	adoptionHandler := NewAdoptionHandler(adoptionService)

	// Root adoptions routes
	router.Post("/", adoptionHandler.CreateAdoption)
	router.Get("/", adoptionHandler.GetAllAdoptions)
	router.Get("/:id", adoptionHandler.GetAdoptionByID)
	router.Put("/:id", adoptionHandler.UpdateAdoption)
	router.Delete("/:id", adoptionHandler.DeleteAdoption)

	// User-specific adoptions routes
	router.Get("/user/:userId", adoptionHandler.GetAdoptionsByUserID)
}
