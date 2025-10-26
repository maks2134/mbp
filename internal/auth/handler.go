package auth

import (
	"errors"
	"mpb/internal/auth/dto"
	"mpb/pkg/errors_constant"
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlers struct {
	AuthService *AuthService
}

func NewAuthHandlers(authService *AuthService) *AuthHandlers {
	return &AuthHandlers{AuthService: authService}
}

func (handler *AuthHandlers) Login(c *fiber.Ctx) error {
	req := middleware.Body[dto.LoginRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	resp, err := handler.AuthService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (handler *AuthHandlers) Register(c *fiber.Ctx) error {
	req := middleware.Body[dto.RegisterRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := handler.AuthService.Register(*req); err != nil {
		if errors.Is(err, errors_constant.UserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}
