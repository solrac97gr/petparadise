package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
)

// PostgresRepository implements the UserRepository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save saves a user into the database
func (r *PostgresRepository) Save(user *models.User) error {
	documentsJSON, err := json.Marshal(user.Documents)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (id, name, email, password, status, created, updated, role, address, phone, documents) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = r.db.Exec(
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Status.String(),
		user.Created,
		user.Updated,
		user.Role.String(),
		user.Address,
		user.Phone,
		documentsJSON,
	)

	return err
}

// FindByID finds a user by their ID
func (r *PostgresRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	var statusStr, roleStr string
	var documentsJSON []byte

	query := `SELECT id, name, email, password, status, created, updated, role, address, phone, documents 
              FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&statusStr,
		&user.Created,
		&user.Updated,
		&roleStr,
		&user.Address,
		&user.Phone,
		&documentsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user.Status = models.Status(statusStr)
	user.Role = models.Role(roleStr)

	if documentsJSON != nil {
		if err := json.Unmarshal(documentsJSON, &user.Documents); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// FindByEmail finds a user by their email
func (r *PostgresRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	var statusStr, roleStr string
	var documentsJSON []byte

	query := `SELECT id, name, email, password, status, created, updated, role, address, phone, documents 
              FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&statusStr,
		&user.Created,
		&user.Updated,
		&roleStr,
		&user.Address,
		&user.Phone,
		&documentsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user.Status = models.Status(statusStr)
	user.Role = models.Role(roleStr)

	if documentsJSON != nil {
		if err := json.Unmarshal(documentsJSON, &user.Documents); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// FindByStatus finds all users with a specific status
func (r *PostgresRepository) FindByStatus(status models.Status) ([]*models.User, error) {
	query := `SELECT id, name, email, password, status, created, updated, role, address, phone, documents 
              FROM users WHERE status = $1`

	rows, err := r.db.Query(query, status.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		var statusStr, roleStr string
		var documentsJSON []byte

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&statusStr,
			&user.Created,
			&user.Updated,
			&roleStr,
			&user.Address,
			&user.Phone,
			&documentsJSON,
		)

		if err != nil {
			return nil, err
		}

		user.Status = models.Status(statusStr)
		user.Role = models.Role(roleStr)

		if documentsJSON != nil {
			if err := json.Unmarshal(documentsJSON, &user.Documents); err != nil {
				return nil, err
			}
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindAll finds all users
func (r *PostgresRepository) FindAll() ([]*models.User, error) {
	query := `SELECT id, name, email, password, status, created, updated, role, address, phone, documents 
              FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		var statusStr, roleStr string
		var documentsJSON []byte

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&statusStr,
			&user.Created,
			&user.Updated,
			&roleStr,
			&user.Address,
			&user.Phone,
			&documentsJSON,
		)

		if err != nil {
			return nil, err
		}

		user.Status = models.Status(statusStr)
		user.Role = models.Role(roleStr)

		if documentsJSON != nil {
			if err := json.Unmarshal(documentsJSON, &user.Documents); err != nil {
				return nil, err
			}
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user
func (r *PostgresRepository) Update(user *models.User) error {
	documentsJSON, err := json.Marshal(user.Documents)
	if err != nil {
		return err
	}

	query := `UPDATE users SET name = $1, email = $2, password = $3, status = $4, updated = $5, 
              role = $6, address = $7, phone = $8, documents = $9 WHERE id = $10`

	_, err = r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Status.String(),
		user.Updated,
		user.Role.String(),
		user.Address,
		user.Phone,
		documentsJSON,
		user.ID,
	)

	return err
}

// Delete deletes a user
func (r *PostgresRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
