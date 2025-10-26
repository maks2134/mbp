package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("body", req)
		return c.Next()
	}
}

func Body[T any](c *fiber.Ctx) *T {
	if v, ok := c.Locals("body").(T); ok {
		return &v
	}
	return nil
}
