package aplication

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/internal/users/domain/ports"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements the UserService interface
type UserService struct {
	repository ports.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repository ports.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, email, password string, role models.Role, address, phone string, documents []string) (*models.User, error) {
	// Check if email is already in use
	existingUser, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	user, err := models.NewUser(
		id,
		name,
		email,
		string(hashedPassword),
		models.StatusActive,
		role,
		address,
		phone,
		documents,
	)
	if err != nil {
		return nil, err
	}

	user.Created = now
	user.Updated = now

	err = s.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID returns a user by their ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repository.FindByID(id)
}

// GetUserByEmail returns a user by their email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repository.FindByEmail(email)
}

// GetUsersByStatus returns all users with a specific status
func (s *UserService) GetUsersByStatus(status models.Status) ([]*models.User, error) {
	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}

	return s.repository.FindByStatus(status)
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repository.FindAll()
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id, name, email, address, phone string, documents []string) (*models.User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// If email is changed, check if the new email is already in use
	if email != user.Email {
		existingUser, err := s.repository.FindByEmail(email)
		if err != nil {
			return nil, err
		}

		if existingUser != nil {
			return nil, errors.New("email already in use")
		}
	}

	if name != "" {
		user.Name = name
	}

	if email != "" {
		user.Email = email
	}

	user.Address = address
	user.Phone = phone

	if documents != nil {
		user.Documents = documents
	}

	user.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(id string, role models.Role) (*models.User, error) {
	if !role.IsValid() {
		return nil, models.ErrInvalidRole
	}

	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Role = role
	user.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserStatus updates a user's status
func (s *UserService) UpdateUserStatus(id string, status models.Status) (*models.User, error) {
	if !status.IsValid() {
		return nil, models.ErrInvalidStatus
	}

	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Status = status
	user.Updated = time.Now().Format(time.RFC3339)

	err = s.repository.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id, oldPassword, newPassword string) error {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Check if old password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Updated = time.Now().Format(time.RFC3339)

	return s.repository.Update(user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id)
}

// Authenticate authenticates a user
func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.Status.IsEquals(models.StatusActive) {
		return nil, errors.New("user account is not active")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
