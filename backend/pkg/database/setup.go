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

	return nil
}
