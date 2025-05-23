package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/adoptions/domain/models"
)

// PostgresRepository implements the AdoptionRepository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save saves an adoption into the database
func (r *PostgresRepository) Save(adoption *models.Adoption) error {
	documentsJSON, err := json.Marshal(adoption.Documents)
	if err != nil {
		return err
	}

	query := `INSERT INTO adoptions (id, pet_id, user_id, status, created, updated, documents) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = r.db.Exec(
		query,
		adoption.ID,
		adoption.PetID,
		adoption.UserID,
		adoption.Status.String(),
		adoption.Created,
		adoption.Updated,
		documentsJSON,
	)

	return err
}

// FindByID finds an adoption by its ID
func (r *PostgresRepository) FindByID(id string) (*models.Adoption, error) {
	var adoption models.Adoption
	var documentsJSON string
	var statusStr string

	query := `SELECT id, pet_id, user_id, status, created, updated, documents 
              FROM adoptions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&adoption.ID,
		&adoption.PetID,
		&adoption.UserID,
		&statusStr,
		&adoption.Created,
		&adoption.Updated,
		&documentsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	adoption.Status = models.Status(statusStr)

	var documents []string
	err = json.Unmarshal([]byte(documentsJSON), &documents)
	if err != nil {
		return nil, err
	}
	
	adoption.Documents = documents

	return &adoption, nil
}

// FindByUserID finds all adoptions for a user
func (r *PostgresRepository) FindByUserID(userID string) ([]*models.Adoption, error) {
	query := `SELECT id, pet_id, user_id, status, created, updated, documents 
              FROM adoptions WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adoptions []*models.Adoption

	for rows.Next() {
		var adoption models.Adoption
		var documentsJSON string
		var statusStr string

		err := rows.Scan(
			&adoption.ID,
			&adoption.PetID,
			&adoption.UserID,
			&statusStr,
			&adoption.Created,
			&adoption.Updated,
			&documentsJSON,
		)

		if err != nil {
			return nil, err
		}

		adoption.Status = models.Status(statusStr)

		var documents []string
		err = json.Unmarshal([]byte(documentsJSON), &documents)
		if err != nil {
			return nil, err
		}
		
		adoption.Documents = documents
		adoptions = append(adoptions, &adoption)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adoptions, nil
}

// FindAll finds all adoptions
func (r *PostgresRepository) FindAll() ([]*models.Adoption, error) {
	query := `SELECT id, pet_id, user_id, status, created, updated, documents FROM adoptions`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adoptions []*models.Adoption

	for rows.Next() {
		var adoption models.Adoption
		var documentsJSON string
		var statusStr string

		err := rows.Scan(
			&adoption.ID,
			&adoption.PetID,
			&adoption.UserID,
			&statusStr,
			&adoption.Created,
			&adoption.Updated,
			&documentsJSON,
		)

		if err != nil {
			return nil, err
		}

		adoption.Status = models.Status(statusStr)

		var documents []string
		err = json.Unmarshal([]byte(documentsJSON), &documents)
		if err != nil {
			return nil, err
		}
		
		adoption.Documents = documents
		adoptions = append(adoptions, &adoption)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adoptions, nil
}

// Update updates an adoption
func (r *PostgresRepository) Update(adoption *models.Adoption) error {
	documentsJSON, err := json.Marshal(adoption.Documents)
	if err != nil {
		return err
	}

	query := `UPDATE adoptions SET pet_id = $1, user_id = $2, status = $3, updated = $4, documents = $5
              WHERE id = $6`
	
	_, err = r.db.Exec(
		query,
		adoption.PetID,
		adoption.UserID,
		adoption.Status.String(),
		adoption.Updated,
		documentsJSON,
		adoption.ID,
	)

	return err
}

// Delete deletes an adoption
func (r *PostgresRepository) Delete(id string) error {
	query := `DELETE FROM adoptions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
