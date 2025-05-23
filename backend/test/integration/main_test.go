package integration

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/solrac97gr/petparadise/pkg/auth"
	"github.com/solrac97gr/petparadise/pkg/config"
	"github.com/solrac97gr/petparadise/pkg/database"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

var testDB *sqlx.DB

func init() {
	// Register command line flags for Godog
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	// Parse command line flags
	flag.Parse()
	opts.Paths = flag.Args()

	// If no paths are specified, use the default features directory
	if len(opts.Paths) == 0 {
		opts.Paths = []string{"features"}
	}

	status := godog.TestSuite{
		Name:                 "Pet Paradise Integration Tests",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options:              &opts,
	}.Run()

	// Exit with non-zero status if there was a failure
	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	// Setup logic for the entire test suite
	// This runs once before any scenarios

	ctx.BeforeSuite(func() {
		// Setup database, perhaps with test data
		fmt.Println("Setting up test suite...")

		// Load configuration for testing
		cfg := config.New()

		// Use test database URL if available, otherwise use default
		dbURL := os.Getenv("TEST_DATABASE_URL")
		if dbURL == "" {
			dbURL = cfg.DatabaseURL
		}

		// Initialize JWT secret
		auth.InitJWTSecret(cfg)

		// Connect to the database
		var err error
		testDB, err = sqlx.Connect("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect to test database: %v", err)
		}

		// Verify database connection
		if err = testDB.Ping(); err != nil {
			log.Fatalf("Failed to ping test database: %v", err)
		}

		// Setup database tables
		if err = database.SetupDatabase(testDB); err != nil {
			log.Fatalf("Failed to setup test database: %v", err)
		}

		fmt.Println("Test suite setup completed successfully")
	})

	ctx.AfterSuite(func() {
		// Cleanup after all scenarios have run
		fmt.Println("Cleaning up test suite...")

		if testDB != nil {
			// Cleanup test data
			testDB.Exec("TRUNCATE TABLE users, pets, adoptions, donations CASCADE")

			// Close database connection
			testDB.Close()
		}

		fmt.Println("Test suite cleanup completed")
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create API client for testing
	apiClient := NewAPIClient()

	// Register step definitions for authentication
	RegisterAuthenticationSteps(ctx, apiClient, testDB)

	// Add hooks for scenario setup/teardown
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Setup before each scenario
		fmt.Printf("Running scenario: %s\n", sc.Name)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// Cleanup after each scenario
		return ctx, nil
	})
}
