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

// Login godoc
// @Summary Login user
// @Description Authenticate user and get JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
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

// Register AuthHandlers godoc
// @Summary Register user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/auth/register [post]
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

// Refresh godoc
// @Summary Refresh tokens
// @Description Get new access and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} dto.RefreshResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *AuthHandlers) Refresh(c *fiber.Ctx) error {
	req := middleware.Body[dto.RefreshRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	resp, err := h.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
