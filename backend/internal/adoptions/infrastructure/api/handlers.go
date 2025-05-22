package api

import (
	"github.com/gofiber/fiber/v2"
)

type AdoptionHandler interface {
	CreateAdoption(c *fiber.Ctx) error
	GetAdoptionByID(c *fiber.Ctx) error
	GetAdoptionsByUserID(c *fiber.Ctx) error
	GetAllAdoptions(c *fiber.Ctx) error
	UpdateAdoption(c *fiber.Ctx) error
	DeleteAdoption(c *fiber.Ctx) error
}
