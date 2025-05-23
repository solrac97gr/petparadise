package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/internal/users/domain/ports"
)

type userHandler struct {
	service ports.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service ports.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

// CreateUser handles the creation of a new user
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	type createUserRequest struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Password  string   `json:"password"`
		Role      string   `json:"role"`
		Address   string   `json:"address"`
		Phone     string   `json:"phone"`
		Documents []string `json:"documents"`
	}

	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	if req.Role == "" {
		req.Role = models.RoleUser.String() // Default role
	}

	role := models.Role(req.Role)
	if !role.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role",
		})
	}

	user, err := h.service.CreateUser(req.Name, req.Email, req.Password, role, req.Address, req.Phone, req.Documents)
	if err != nil {
		if err.Error() == "email already in use" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return the password, even though it's already marked as json:"-"
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID handles getting a single user by ID
func (h *userHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(user)
}

// GetUserByEmail handles getting a user by email
func (h *userHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(user)
}

// GetUsersByStatus handles getting all users with a specific status
func (h *userHandler) GetUsersByStatus(c *fiber.Ctx) error {
	statusParam := c.Query("status")
	if statusParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status query parameter is required",
		})
	}

	status := models.Status(statusParam)
	if !status.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	users, err := h.service.GetUsersByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return passwords
	for _, user := range users {
		user.Password = ""
	}

	return c.JSON(users)
}

// GetAllUsers handles getting all users
func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return passwords
	for _, user := range users {
		user.Password = ""
	}

	return c.JSON(users)
}

// UpdateUser handles updating a user's information
func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type updateUserRequest struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Address   string   `json:"address"`
		Phone     string   `json:"phone"`
		Documents []string `json:"documents"`
	}

	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.service.UpdateUser(id, req.Name, req.Email, req.Address, req.Phone, req.Documents)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err.Error() == "email already in use" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(user)
}

// UpdateUserRole handles updating a user's role
func (h *userHandler) UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type updateRoleRequest struct {
		Role string `json:"role"`
	}

	var req updateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role is required",
		})
	}

	role := models.Role(req.Role)
	if !role.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role",
		})
	}

	user, err := h.service.UpdateUserRole(id, role)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(user)
}

// UpdateUserStatus handles updating a user's status
func (h *userHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type updateStatusRequest struct {
		Status string `json:"status"`
	}

	var req updateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	status := models.Status(req.Status)
	if !status.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	user, err := h.service.UpdateUserStatus(id, status)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(user)
}

// ChangePassword handles changing a user's password
func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type changePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var req changePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.OldPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Old password is required",
		})
	}

	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password is required",
		})
	}

	err := h.service.ChangePassword(id, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err.Error() == "incorrect password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// DeleteUser handles deleting a user
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	err := h.service.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// Login handles user authentication
func (h *userHandler) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Don't return the password
	user.Password = ""

	return c.JSON(fiber.Map{
		"user": user,
		// We would normally generate and return a JWT token here
		"token": "sample-jwt-token", // placeholder - in a real app, this would be a real JWT
	})
}

// Logout handles user logout
func (h *userHandler) Logout(c *fiber.Ctx) error {
	// In a real app, we would invalidate the JWT token
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
