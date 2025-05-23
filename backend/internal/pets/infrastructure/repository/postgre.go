package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/pets/domain/models"
)

// PostgresRepository implements the PetRepository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save saves a pet into the database
func (r *PostgresRepository) Save(pet *models.Pet) error {
	imagesJSON, err := json.Marshal(pet.Images)
	if err != nil {
		return err
	}

	query := `INSERT INTO pets (id, name, species, breed, age, description, status, created, updated, images) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = r.db.Exec(
		query,
		pet.ID,
		pet.Name,
		pet.Species,
		pet.Breed,
		pet.Age,
		pet.Description,
		pet.Status.String(),
		pet.Created,
		pet.Updated,
		imagesJSON,
	)

	return err
}

// FindByID finds a pet by its ID
func (r *PostgresRepository) FindByID(id string) (*models.Pet, error) {
	var pet models.Pet
	var imagesJSON string
	var statusStr string

	query := `SELECT id, name, species, breed, age, description, status, created, updated, images 
              FROM pets WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&pet.ID,
		&pet.Name,
		&pet.Species,
		&pet.Breed,
		&pet.Age,
		&pet.Description,
		&statusStr,
		&pet.Created,
		&pet.Updated,
		&imagesJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	pet.Status = models.Status(statusStr)

	var images []string
	err = json.Unmarshal([]byte(imagesJSON), &images)
	if err != nil {
		return nil, err
	}

	pet.Images = images

	return &pet, nil
}

// FindByStatus finds all pets with a specific status
func (r *PostgresRepository) FindByStatus(status models.Status) ([]*models.Pet, error) {
	query := `SELECT id, name, species, breed, age, description, status, created, updated, images 
              FROM pets WHERE status = $1`

	rows, err := r.db.Query(query, status.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []*models.Pet

	for rows.Next() {
		var pet models.Pet
		var imagesJSON string
		var statusStr string

		err := rows.Scan(
			&pet.ID,
			&pet.Name,
			&pet.Species,
			&pet.Breed,
			&pet.Age,
			&pet.Description,
			&statusStr,
			&pet.Created,
			&pet.Updated,
			&imagesJSON,
		)

		if err != nil {
			return nil, err
		}

		pet.Status = models.Status(statusStr)

		var images []string
		err = json.Unmarshal([]byte(imagesJSON), &images)
		if err != nil {
			return nil, err
		}

		pet.Images = images
		pets = append(pets, &pet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pets, nil
}

// FindAll finds all pets
func (r *PostgresRepository) FindAll() ([]*models.Pet, error) {
	query := `SELECT id, name, species, breed, age, description, status, created, updated, images FROM pets`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []*models.Pet

	for rows.Next() {
		var pet models.Pet
		var imagesJSON string
		var statusStr string

		err := rows.Scan(
			&pet.ID,
			&pet.Name,
			&pet.Species,
			&pet.Breed,
			&pet.Age,
			&pet.Description,
			&statusStr,
			&pet.Created,
			&pet.Updated,
			&imagesJSON,
		)

		if err != nil {
			return nil, err
		}

		pet.Status = models.Status(statusStr)

		var images []string
		err = json.Unmarshal([]byte(imagesJSON), &images)
		if err != nil {
			return nil, err
		}

		pet.Images = images
		pets = append(pets, &pet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pets, nil
}

// Update updates a pet
func (r *PostgresRepository) Update(pet *models.Pet) error {
	imagesJSON, err := json.Marshal(pet.Images)
	if err != nil {
		return err
	}

	query := `UPDATE pets SET name = $1, species = $2, breed = $3, age = $4, description = $5, 
              status = $6, updated = $7, images = $8 WHERE id = $9`

	_, err = r.db.Exec(
		query,
		pet.Name,
		pet.Species,
		pet.Breed,
		pet.Age,
		pet.Description,
		pet.Status.String(),
		pet.Updated,
		imagesJSON,
		pet.ID,
	)

	return err
}

// Delete deletes a pet
func (r *PostgresRepository) Delete(id string) error {
	query := `DELETE FROM pets WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
