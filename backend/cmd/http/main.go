package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/solrac97gr/petparadise/pkg/config"
	"github.com/solrac97gr/petparadise/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize logger
	appLogger := logger.New(cfg.LogLevel)
	defer appLogger.Sync()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Pet Paradise API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Pet Paradise API!")
	})

	// Health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes - would be connected to the proper handlers
	api := app.Group("/api")

	// Users routes
	users := api.Group("/users")
	users.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get all users endpoint"})
	})

	// Pets routes
	pets := api.Group("/pets")
	pets.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get all pets endpoint"})
	})

	// Adoptions routes
	adoptions := api.Group("/adoptions")
	adoptions.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get all adoptions endpoint"})
	})

	// Donations routes
	donations := api.Group("/donations")
	donations.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get all donations endpoint"})
	})

	// Start server
	serverPort := strconv.Itoa(cfg.ServerPort)
	appLogger.Info("Starting server on port " + serverPort)
	log.Fatal(app.Listen(":" + serverPort))
}
