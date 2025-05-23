package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/petparadise/internal/adoptions/domain/models"
	"github.com/solrac97gr/petparadise/internal/adoptions/domain/ports"
)

type adoptionHandler struct {
	service ports.AdoptionService
}

// NewAdoptionHandler creates a new adoption handler
func NewAdoptionHandler(service ports.AdoptionService) AdoptionHandler {
	return &adoptionHandler{
		service: service,
	}
}

// CreateAdoption handles the creation of a new adoption
func (h *adoptionHandler) CreateAdoption(c *fiber.Ctx) error {
	type createAdoptionRequest struct {
		PetID     string   `json:"pet_id"`
		UserID    string   `json:"user_id"`
		Documents []string `json:"documents"`
	}

	var req createAdoptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.PetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pet ID is required",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	adoption, err := h.service.CreateAdoption(req.PetID, req.UserID, req.Documents)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(adoption)
}

// GetAdoptionByID handles getting a single adoption by ID
func (h *adoptionHandler) GetAdoptionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	adoption, err := h.service.GetAdoptionByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if adoption == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Adoption not found",
		})
	}

	return c.JSON(adoption)
}

// GetAdoptionsByUserID handles getting all adoptions for a specific user
func (h *adoptionHandler) GetAdoptionsByUserID(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	adoptions, err := h.service.GetAdoptionsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(adoptions)
}

// GetAllAdoptions handles getting all adoptions
func (h *adoptionHandler) GetAllAdoptions(c *fiber.Ctx) error {
	adoptions, err := h.service.GetAllAdoptions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(adoptions)
}

// UpdateAdoption handles updating an adoption
func (h *adoptionHandler) UpdateAdoption(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type updateAdoptionRequest struct {
		Status    string   `json:"status"`
		Documents []string `json:"documents,omitempty"`
	}

	var req updateAdoptionRequest
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

	adoption, err := h.service.UpdateAdoption(id, status, req.Documents)
	if err != nil {
		if err == models.ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if adoption == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Adoption not found",
		})
	}

	return c.JSON(adoption)
}

// DeleteAdoption handles deleting an adoption
func (h *adoptionHandler) DeleteAdoption(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	err := h.service.DeleteAdoption(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
