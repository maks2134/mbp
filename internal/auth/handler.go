package auth

import "github.com/gofiber/fiber/v2"

type AuthHandlers struct {
	AuthService *AuthService
}

func NewAuthHandlers(authService *AuthService) *AuthHandlers {
	return &AuthHandlers{AuthService: authService}
}

func (handler *AuthHandlers) Login(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{})
}

func (handler *AuthHandlers) Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
