package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
)

// Protected is a middleware that checks if the request has a valid JWT token
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check if the header is in the correct format (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract token from header
		tokenString := parts[1]

		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			status := fiber.StatusUnauthorized
			errorMsg := "Invalid token"

			if err == ErrExpiredToken {
				errorMsg = "Token has expired"
			}

			return c.Status(status).JSON(fiber.Map{
				"error": errorMsg,
			})
		}

		// Set claims in context for use in handlers
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		// Continue to next middleware/handler
		return c.Next()
	}
}

// RoleRequired is a middleware that checks if the user has the required role
func RoleRequired(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by Protected middleware)
		userRole, ok := c.Locals("role").(models.Role)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Failed to get user role",
			})
		}

		// Admin role always has access
		if userRole.IsEquals(models.RoleAdmin) {
			return c.Next()
		}

		// Check if user has one of the required roles
		for _, role := range roles {
			if userRole.IsEquals(role) {
				return c.Next()
			}
		}

		// If we get here, the user doesn't have the required role
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
