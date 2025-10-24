package auth

import "github.com/gofiber/fiber/v2"

type AuthHandlers struct {
	AuthRepository *AuthRepository
}

func NewAuthHandlers(authRepository *AuthRepository) *AuthHandlers {
	return &AuthHandlers{AuthRepository: authRepository}
}

func (handler *AuthHandlers) Login(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (handler *AuthHandlers) Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
