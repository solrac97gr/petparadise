package config

import (
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	ServerPort         int
	DatabaseURL        string
	JWTSecret          string
	LogLevel           string
	Environment        string
	CORSAllowedOrigins []string
}

// New creates a new configuration instance with values from environment variables
func New() *Config {
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "3000"))

	return &Config{
		ServerPort:         port,
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/petparadise?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		CORSAllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")},
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
