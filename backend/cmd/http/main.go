package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	adoptionAPI "github.com/solrac97gr/petparadise/internal/adoptions/infrastructure/api"
	donationAPI "github.com/solrac97gr/petparadise/internal/donations/infrastructure/api"
	petAPI "github.com/solrac97gr/petparadise/internal/pets/infrastructure/api"
	userAPI "github.com/solrac97gr/petparadise/internal/users/infrastructure/api"
	"github.com/solrac97gr/petparadise/pkg/auth"
	"github.com/solrac97gr/petparadise/pkg/config"
	"github.com/solrac97gr/petparadise/pkg/database"
	"github.com/solrac97gr/petparadise/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize logger
	appLogger := logger.New(cfg.LogLevel)
	defer appLogger.Sync()

	// Initialize JWT secret
	auth.InitJWTSecret(cfg)

	// Connect to the database
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	// Verify database connection
	if err = db.Ping(); err != nil {
		appLogger.Fatal("Failed to ping database: " + err.Error())
	}

	// Setup database tables
	if err = database.SetupDatabase(db); err != nil {
		appLogger.Fatal("Failed to setup database: " + err.Error())
	}

	appLogger.Info("Connected to database")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Pet Paradise API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSAllowedOrigins, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH",
		AllowCredentials: true,
		ExposeHeaders:    "Authorization",
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

	// API routes
	api := app.Group("/api")

	// Users routes
	users := api.Group("/users")
	userAPI.SetupUserRoutes(users, db)

	// Pets routes
	pets := api.Group("/pets")
	petAPI.SetupPetRoutes(pets, db)

	// Adoptions routes
	adoptions := api.Group("/adoptions")
	adoptionAPI.SetupAdoptionRoutes(adoptions, db)

	// Donations routes
	donations := api.Group("/donations")
	donationAPI.SetupDonationRoutes(donations, db)

	// Start server
	serverPort := strconv.Itoa(cfg.ServerPort)
	appLogger.Info("Starting server on port " + serverPort)
	log.Fatal(app.Listen(":" + serverPort))
}
