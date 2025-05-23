package database

import (
	"github.com/jmoiron/sqlx"
)

// SetupDatabase initializes database tables if they don't exist
func SetupDatabase(db *sqlx.DB) error {
	// Create adoptions table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS adoptions (
			id VARCHAR(36) PRIMARY KEY,
			pet_id VARCHAR(36) NOT NULL,
			user_id VARCHAR(36) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created TIMESTAMP NOT NULL,
			updated TIMESTAMP NOT NULL,
			documents JSONB
		);
		
		CREATE INDEX IF NOT EXISTS idx_adoptions_pet_id ON adoptions(pet_id);
		CREATE INDEX IF NOT EXISTS idx_adoptions_user_id ON adoptions(user_id);
		CREATE INDEX IF NOT EXISTS idx_adoptions_status ON adoptions(status);
	`)
	if err != nil {
		return err
	}

	// Create pets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pets (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			species VARCHAR(50) NOT NULL,
			breed VARCHAR(100),
			age INT NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL,
			created TIMESTAMP NOT NULL,
			updated TIMESTAMP NOT NULL,
			images JSONB
		);
		
		CREATE INDEX IF NOT EXISTS idx_pets_status ON pets(status);
		CREATE INDEX IF NOT EXISTS idx_pets_species ON pets(species);
		CREATE INDEX IF NOT EXISTS idx_pets_breed ON pets(breed);
	`)
	if err != nil {
		return err
	}

	// Create donations table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created TIMESTAMP NOT NULL,
			updated TIMESTAMP NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
	`)
	if err != nil {
		return err
	}

	// Create donations table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS donations (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created TIMESTAMP NOT NULL,
			updated TIMESTAMP NOT NULL,
			comment TEXT,
			anonymous BOOLEAN NOT NULL DEFAULT FALSE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		
		CREATE INDEX IF NOT EXISTS idx_donations_user_id ON donations(user_id);
		CREATE INDEX IF NOT EXISTS idx_donations_status ON donations(status);
		CREATE INDEX IF NOT EXISTS idx_donations_created ON donations(created);
	`)
	if err != nil {
		return err
	}

	return nil
}
