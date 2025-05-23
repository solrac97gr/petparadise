package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solrac97gr/petparadise/internal/donations/domain/models"
	"github.com/solrac97gr/petparadise/internal/donations/domain/ports"
)

type donationHandler struct {
	service ports.DonationService
}

// NewDonationHandler creates a new donation handler
func NewDonationHandler(service ports.DonationService) DonationHandler {
	return &donationHandler{
		service: service,
	}
}

// CreateDonation handles the creation of a new donation
func (h *donationHandler) CreateDonation(c *fiber.Ctx) error {
	type createDonationRequest struct {
		UserID    string  `json:"user_id"`
		Amount    float64 `json:"amount"`
		Comment   string  `json:"comment"`
		Anonymous bool    `json:"anonymous"`
	}

	var req createDonationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID is required",
		})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	donation, err := h.service.CreateDonation(req.UserID, req.Amount, req.Comment, req.Anonymous)
	if err != nil {
		if err == models.ErrInvalidAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(donation)
}

// GetDonationByID handles getting a single donation by ID
func (h *donationHandler) GetDonationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	donation, err := h.service.GetDonationByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if donation == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Donation not found",
		})
	}

	return c.JSON(donation)
}

// GetDonationsByUserID handles getting all donations for a user
func (h *donationHandler) GetDonationsByUserID(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	donations, err := h.service.GetDonationsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(donations)
}

// GetAllDonations handles getting all donations
func (h *donationHandler) GetAllDonations(c *fiber.Ctx) error {
	donations, err := h.service.GetAllDonations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(donations)
}

// UpdateDonationStatus handles updating a donation's status
func (h *donationHandler) UpdateDonationStatus(c *fiber.Ctx) error {
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

	donation, err := h.service.UpdateDonation(id, status)
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

	if donation == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Donation not found",
		})
	}

	return c.JSON(donation)
}

// DeleteDonation handles deleting a donation
func (h *donationHandler) DeleteDonation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	err := h.service.DeleteDonation(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
