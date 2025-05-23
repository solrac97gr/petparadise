package integration

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

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
		// connect to test database
		// migrate tables
		// seed initial data
	})

	ctx.AfterSuite(func() {
		// Cleanup after all scenarios have run
		fmt.Println("Cleaning up test suite...")
		// close database connections
		// remove test data
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create API client for testing
	apiClient := NewAPIClient()

	// Register step definitions for authentication
	RegisterAuthenticationSteps(ctx, apiClient)

	// Register step definitions for users
	RegisterUserSteps(ctx, apiClient)

	// Register step definitions for pets
	RegisterPetSteps(ctx, apiClient)

	// Register step definitions for adoptions
	RegisterAdoptionSteps(ctx, apiClient)

	// Register step definitions for donations
	RegisterDonationSteps(ctx, apiClient)

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
