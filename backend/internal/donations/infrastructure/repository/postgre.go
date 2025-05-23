package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/donations/domain/models"
)

// PostgresRepository implements the DonationRepository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save saves a donation into the database
func (r *PostgresRepository) Save(donation *models.Donation) error {
	query := `INSERT INTO donations (id, user_id, amount, status, created, updated, comment, anonymous) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(
		query,
		donation.ID,
		donation.UserID,
		donation.Amount,
		donation.Status.String(),
		donation.Created,
		donation.Updated,
		donation.Comment,
		donation.Anonymous,
	)

	return err
}

// FindByID finds a donation by its ID
func (r *PostgresRepository) FindByID(id string) (*models.Donation, error) {
	var donation models.Donation
	var statusStr string

	query := `SELECT id, user_id, amount, status, created, updated, comment, anonymous 
              FROM donations WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&donation.ID,
		&donation.UserID,
		&donation.Amount,
		&statusStr,
		&donation.Created,
		&donation.Updated,
		&donation.Comment,
		&donation.Anonymous,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	donation.Status = models.Status(statusStr)

	return &donation, nil
}

// FindByUserID finds all donations for a user
func (r *PostgresRepository) FindByUserID(userID string) ([]*models.Donation, error) {
	query := `SELECT id, user_id, amount, status, created, updated, comment, anonymous 
              FROM donations WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []*models.Donation

	for rows.Next() {
		var donation models.Donation
		var statusStr string

		err := rows.Scan(
			&donation.ID,
			&donation.UserID,
			&donation.Amount,
			&statusStr,
			&donation.Created,
			&donation.Updated,
			&donation.Comment,
			&donation.Anonymous,
		)

		if err != nil {
			return nil, err
		}

		donation.Status = models.Status(statusStr)
		donations = append(donations, &donation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return donations, nil
}

// FindAll finds all donations
func (r *PostgresRepository) FindAll() ([]*models.Donation, error) {
	query := `SELECT id, user_id, amount, status, created, updated, comment, anonymous FROM donations`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []*models.Donation

	for rows.Next() {
		var donation models.Donation
		var statusStr string

		err := rows.Scan(
			&donation.ID,
			&donation.UserID,
			&donation.Amount,
			&statusStr,
			&donation.Created,
			&donation.Updated,
			&donation.Comment,
			&donation.Anonymous,
		)

		if err != nil {
			return nil, err
		}

		donation.Status = models.Status(statusStr)
		donations = append(donations, &donation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return donations, nil
}

// Update updates a donation
func (r *PostgresRepository) Update(donation *models.Donation) error {
	query := `UPDATE donations SET user_id = $1, amount = $2, status = $3, updated = $4, 
              comment = $5, anonymous = $6 WHERE id = $7`

	_, err := r.db.Exec(
		query,
		donation.UserID,
		donation.Amount,
		donation.Status.String(),
		donation.Updated,
		donation.Comment,
		donation.Anonymous,
		donation.ID,
	)

	return err
}

// Delete deletes a donation
func (r *PostgresRepository) Delete(id string) error {
	query := `DELETE FROM donations WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
