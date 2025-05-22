package api

import (
	"github.com/gofiber/fiber/v2"
)

type PetHandler interface {
	CreatePet(c *fiber.Ctx) error
	GetPetByID(c *fiber.Ctx) error
	GetPetsByStatus(c *fiber.Ctx) error
	GetAllPets(c *fiber.Ctx) error
	UpdatePet(c *fiber.Ctx) error
	UpdatePetStatus(c *fiber.Ctx) error
	DeletePet(c *fiber.Ctx) error
}
