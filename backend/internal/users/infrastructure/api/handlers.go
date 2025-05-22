package api

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetUserByEmail(c *fiber.Ctx) error
	GetUsersByStatus(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	UpdateUserRole(c *fiber.Ctx) error
	UpdateUserStatus(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}
