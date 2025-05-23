package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/users/aplication"
	"github.com/solrac97gr/petparadise/internal/users/infrastructure/repository"
)

// SetupUserRoutes sets up all user routes
func SetupUserRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	userRepo := repository.NewPostgresRepository(db)

	// Initialize service
	userService := aplication.NewUserService(userRepo)

	// Initialize handler
	userHandler := NewUserHandler(userService)

	// User routes
	router.Post("/", userHandler.CreateUser)
	router.Get("/", userHandler.GetAllUsers)
	router.Get("/:id", userHandler.GetUserByID)
	router.Get("/email", userHandler.GetUserByEmail)
	router.Get("/status", userHandler.GetUsersByStatus)
	router.Put("/:id", userHandler.UpdateUser)
	router.Patch("/:id/role", userHandler.UpdateUserRole)
	router.Patch("/:id/status", userHandler.UpdateUserStatus)
	router.Post("/:id/password", userHandler.ChangePassword)
	router.Delete("/:id", userHandler.DeleteUser)

	// Auth routes
	router.Post("/login", userHandler.Login)
	router.Post("/logout", userHandler.Logout)
}
