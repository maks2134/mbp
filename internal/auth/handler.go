package auth

import "github.com/gofiber/fiber/v2"

type AuthHandlers struct {
	AuthRepository *AuthRepository
}

func (handler *AuthHandlers) Login(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{})
}
