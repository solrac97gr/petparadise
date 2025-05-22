package api

import (
	"github.com/gofiber/fiber/v2"
)

type DonationHandler interface {
	CreateDonation(c *fiber.Ctx) error
	GetDonationByID(c *fiber.Ctx) error
	GetDonationsByUserID(c *fiber.Ctx) error
	GetAllDonations(c *fiber.Ctx) error
	UpdateDonation(c *fiber.Ctx) error
	DeleteDonation(c *fiber.Ctx) error
}
