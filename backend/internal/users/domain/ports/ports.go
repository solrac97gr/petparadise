package ports

import "github.com/solrac97gr/petparadise/internal/users/domain/models"

type UserRepository interface {
	Save(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByStatus(status models.Status) ([]*models.User, error)
	FindAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type UserService interface {
	CreateUser(name, email, password string, role models.Role, address, phone string, documents []string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsersByStatus(status models.Status) ([]*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id, name, email, address, phone string, documents []string) (*models.User, error)
	UpdateUserRole(id string, role models.Role) (*models.User, error)
	UpdateUserStatus(id string, status models.Status) (*models.User, error)
	ChangePassword(id, oldPassword, newPassword string) error
	DeleteUser(id string) error
	Authenticate(email, password string) (*models.User, error)
}
