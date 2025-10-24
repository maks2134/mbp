package auth

import "github.com/gofiber/fiber/v2"

type AuthRoutes struct {
	router fiber.Router
}

func NewAuthRoutes(router fiber.Router) *AuthRoutes {
	return &AuthRoutes{router: router}
}

func (r *AuthRoutes) Register() {
	auth := r.router.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"msg": "login ok"}) })
	auth.Post("/register", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"msg": "register ok"}) })
	auth.Post("/logout", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"msg": "logout ok"}) })
}
