package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/petparadise/internal/pets/domain/models"
	"github.com/solrac97gr/petparadise/internal/pets/domain/ports"
)

type petHandler struct {
	service ports.PetService
}

// NewPetHandler creates a new pet handler
func NewPetHandler(service ports.PetService) PetHandler {
	return &petHandler{
		service: service,
	}
}

// CreatePet handles the creation of a new pet
func (h *petHandler) CreatePet(c *fiber.Ctx) error {
	type createPetRequest struct {
		Name        string   `json:"name"`
		Species     string   `json:"species"`
		Breed       string   `json:"breed"`
		Age         int      `json:"age"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	}

	var req createPetRequest
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

	if req.Species == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Species is required",
		})
	}

	if req.Age < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Age must be a positive number",
		})
	}

	pet, err := h.service.CreatePet(req.Name, req.Species, req.Breed, req.Age, req.Description, req.Images)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(pet)
}

// GetPetByID handles getting a single pet by ID
func (h *petHandler) GetPetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	pet, err := h.service.GetPetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if pet == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pet not found",
		})
	}

	return c.JSON(pet)
}

// GetPetsByStatus handles getting all pets with a specific status
func (h *petHandler) GetPetsByStatus(c *fiber.Ctx) error {
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

	pets, err := h.service.GetPetsByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(pets)
}

// GetAllPets handles getting all pets
func (h *petHandler) GetAllPets(c *fiber.Ctx) error {
	pets, err := h.service.GetAllPets()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(pets)
}

// UpdatePet handles updating a pet
func (h *petHandler) UpdatePet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	type updatePetRequest struct {
		Name        string   `json:"name"`
		Species     string   `json:"species"`
		Breed       string   `json:"breed"`
		Age         *int     `json:"age"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		Images      []string `json:"images"`
	}

	var req updatePetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	age := 0
	if req.Age != nil {
		age = *req.Age
	}

	var status models.Status
	if req.Status != "" {
		status = models.Status(req.Status)
		if !status.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status",
			})
		}
	}

	pet, err := h.service.UpdatePet(id, req.Name, req.Species, req.Breed, age, req.Description, status, req.Images)
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

	if pet == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pet not found",
		})
	}

	return c.JSON(pet)
}

// UpdatePetStatus handles updating only a pet's status
func (h *petHandler) UpdatePetStatus(c *fiber.Ctx) error {
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

	pet, err := h.service.UpdatePetStatus(id, status)
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

	if pet == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pet not found",
		})
	}

	return c.JSON(pet)
}

// DeletePet handles deleting a pet
func (h *petHandler) DeletePet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	err := h.service.DeletePet(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
