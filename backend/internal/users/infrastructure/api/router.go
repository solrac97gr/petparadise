package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/users/aplication"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/internal/users/infrastructure/repository"
	"github.com/solrac97gr/petparadise/pkg/auth"
)

// SetupUserRoutes sets up all user routes
func SetupUserRoutes(router fiber.Router, db *sqlx.DB) {
	// Initialize repository
	userRepo := repository.NewPostgresRepository(db)

	// Initialize service
	userService := aplication.NewUserService(userRepo)

	// Initialize handler
	userHandler := NewUserHandler(userService)

	// Public routes
	router.Post("/register", userHandler.CreateUser)  // Registration endpoint
	router.Post("/login", userHandler.Login)          // Login endpoint
	router.Post("/refresh", userHandler.RefreshToken) // Token refresh endpoint

	// Protected routes
	protectedRoutes := router.Use(auth.Protected())

	// Specific routes first (before parameterized routes)
	// Logout route (protected)
	protectedRoutes.Post("/logout", userHandler.Logout)

	// User management routes (protected)
	protectedRoutes.Get("/", userHandler.GetAllUsers)
	protectedRoutes.Get("/:id", userHandler.GetUserByID)
	protectedRoutes.Get("/email", userHandler.GetUserByEmail)
	protectedRoutes.Get("/status", userHandler.GetUsersByStatus)
	protectedRoutes.Put("/:id", userHandler.UpdateUser)

	// Admin only routes
	adminRoutes := protectedRoutes.Use(auth.RoleRequired(models.RoleAdmin))
	adminRoutes.Patch("/:id/role", userHandler.UpdateUserRole)
	adminRoutes.Patch("/:id/status", userHandler.UpdateUserStatus)
	adminRoutes.Delete("/:id", userHandler.DeleteUser)

	// User password management (protected)
	protectedRoutes.Post("/:id/password", userHandler.ChangePassword)

	// Revoke all tokens for a user (protected)
	protectedRoutes.Post("/:id/revoke-tokens", userHandler.RevokeUserTokens)
}
