package auth

import "github.com/gofiber/fiber/v2"

type AuthRoutes struct {
	router      fiber.Router
	authHandler *AuthHandlers
}

func NewAuthRoutes(router fiber.Router, authHandler *AuthHandlers) *AuthRoutes {
	return &AuthRoutes{router: router, authHandler: authHandler}
}

func (r *AuthRoutes) Register() {
	auth := r.router.Group("/auth")
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/register", r.authHandler.Login)
	auth.Post("/logout", r.authHandler.Login)
}
