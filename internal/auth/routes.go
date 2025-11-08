package auth

import (
	"mpb/internal/auth/dto"
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthRoutes struct {
	router  fiber.Router
	handler *AuthHandlers
}

func NewAuthRoutes(router fiber.Router, handler *AuthHandlers) *AuthRoutes {
	return &AuthRoutes{router: router, handler: handler}
}

func (r *AuthRoutes) Register() {
	auth := r.router.Group("/auth")

	auth.Post("/register",
		middleware.ValidateBody[dto.RegisterRequest](),
		r.handler.Register,
	)

	auth.Post("/login",
		middleware.ValidateBody[dto.LoginRequest](),
		r.handler.Login,
	)

	auth.Post("/refresh",
		middleware.ValidateBody[dto.RefreshRequest](),
		r.handler.Refresh,
	)
}
