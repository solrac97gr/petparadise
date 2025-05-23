package integration

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SetupTestData creates test users and other necessary data for testing
func SetupTestData(db *sqlx.DB) error {
	fmt.Println("Starting test data setup...")

	// Create test users for different roles
	testUsers := []struct {
		ID       string
		Name     string
		Email    string
		Password string
		Role     string
		Status   string
	}{
		{
			ID:       "11111111-1111-1111-1111-111111111111",
			Name:     "Test Admin",
			Email:    "admin@example.com",
			Password: "password123",
			Role:     "admin",
			Status:   "active",
		},
		{
			ID:       "22222222-2222-2222-2222-222222222222",
			Name:     "Test User",
			Email:    "user@example.com",
			Password: "password123",
			Role:     "user",
			Status:   "active",
		},
		{
			ID:       "33333333-3333-3333-3333-333333333333",
			Name:     "Test Volunteer",
			Email:    "volunteer@example.com",
			Password: "password123",
			Role:     "volunteer",
			Status:   "active",
		},
		{
			ID:       "44444444-4444-4444-4444-444444444444",
			Name:     "Test Vet",
			Email:    "vet@example.com",
			Password: "password123",
			Role:     "vet",
			Status:   "active",
		},
		{
			ID:       "55555555-5555-5555-5555-555555555555",
			Name:     "John Doe",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "user",
			Status:   "active",
		},
	}

	// Insert test users
	for _, user := range testUsers {
		fmt.Printf("Setting up user: %s (%s)\n", user.Name, user.Email)

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error hashing password for %s: %v\n", user.Email, err)
			return err
		}

		// Check if user already exists
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", user.Email)
		if err != nil {
			fmt.Printf("Error checking user existence for %s: %v\n", user.Email, err)
			return err
		}

		if count == 0 {
			fmt.Printf("Inserting user: %s\n", user.Email)
			// Insert user
			_, err = db.Exec(`
				INSERT INTO users (id, name, email, password, status, role, created, updated, address, phone, documents)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), '', '', '[]'::jsonb)
				ON CONFLICT (email) DO NOTHING
			`, user.ID, user.Name, user.Email, string(hashedPassword), user.Status, user.Role)
			if err != nil {
				fmt.Printf("Error inserting user %s: %v\n", user.Email, err)
				return err
			}
			fmt.Printf("Successfully inserted user: %s\n", user.Email)
		} else {
			fmt.Printf("User %s already exists, skipping\n", user.Email)
		}
	}

	// Create test pets
	testPets := []struct {
		ID          string
		Name        string
		Species     string
		Breed       string
		Age         int
		Description string
		Status      string
		Images      string
	}{
		{
			ID:          "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			Name:        "Fluffy",
			Species:     "Dog",
			Breed:       "Golden Retriever",
			Age:         3,
			Description: "A friendly and energetic dog",
			Status:      "available",
			Images:      "[]",
		},
		{
			ID:          "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			Name:        "Whiskers",
			Species:     "Cat",
			Breed:       "Persian",
			Age:         2,
			Description: "A calm and gentle cat",
			Status:      "available",
			Images:      "[]",
		},
	}

	// Insert test pets
	for _, pet := range testPets {
		// Check if pet already exists
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM pets WHERE id = $1", pet.ID)
		if err != nil {
			return err
		}

		if count == 0 {
			// Insert pet
			_, err = db.Exec(`
				INSERT INTO pets (id, name, species, breed, age, description, status, created, updated, images)
				VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), $8::jsonb)
				ON CONFLICT (id) DO NOTHING
			`, pet.ID, pet.Name, pet.Species, pet.Breed, pet.Age, pet.Description, pet.Status, pet.Images)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanupTestData removes test data after tests complete
func CleanupTestData(db *sqlx.DB) error {
	// Remove test users
	testEmails := []string{
		"admin@example.com",
		"user@example.com",
		"volunteer@example.com",
		"vet@example.com",
		"test@example.com",
		"john@example.com",
	}

	for _, email := range testEmails {
		_, err := db.Exec("DELETE FROM users WHERE email = $1", email)
		if err != nil {
			return err
		}
	}

	// Remove test pets
	testPetIDs := []string{
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
	}
	for _, id := range testPetIDs {
		_, err := db.Exec("DELETE FROM pets WHERE id = $1", id)
		if err != nil {
			return err
		}
	}

	// Remove test adoptions
	_, err := db.Exec("DELETE FROM adoptions WHERE pet_id IN ($1, $2)", testPetIDs[0], testPetIDs[1])
	if err != nil {
		return err
	}

	// Remove test donations using proper UUIDs
	testUserIDs := []string{
		"11111111-1111-1111-1111-111111111111", // admin
		"22222222-2222-2222-2222-222222222222", // user
		"33333333-3333-3333-3333-333333333333", // volunteer
		"44444444-4444-4444-4444-444444444444", // vet
	}
	for _, userID := range testUserIDs {
		_, err = db.Exec("DELETE FROM donations WHERE user_id = $1", userID)
		if err != nil {
			return err
		}
	}

	return nil
}
