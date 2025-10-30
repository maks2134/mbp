package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T

		if err := c.BodyParser(&req); err != nil {
			return badRequest(c, "invalid request body")
		}

		if err := validate.Struct(req); err != nil {
			return validationError(c, err)
		}

		c.Locals("body", req)
		return c.Next()
	}
}

func ValidateQuery[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T

		if err := c.QueryParser(&req); err != nil {
			return badRequest(c, "invalid query parameters")
		}

		if err := validate.Struct(req); err != nil {
			return validationError(c, err)
		}

		c.Locals("query", req)
		return c.Next()
	}
}

func ValidateHeader[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T

		if err := c.ReqHeaderParser(&req); err != nil {
			return badRequest(c, "invalid headers")
		}

		if err := validate.Struct(req); err != nil {
			return validationError(c, err)
		}

		c.Locals("header", req)
		return c.Next()
	}
}

func Body[T any](c *fiber.Ctx) *T {
	if v := c.Locals("body"); v != nil {
		if typed, ok := v.(T); ok {
			return &typed
		}
	}
	return nil
}

func Query[T any](c *fiber.Ctx) *T {
	if v := c.Locals("query"); v != nil {
		if typed, ok := v.(T); ok {
			return &typed
		}
	}
	return nil
}

func Header[T any](c *fiber.Ctx) *T {
	if v := c.Locals("header"); v != nil {
		if typed, ok := v.(T); ok {
			return &typed
		}
	}
	return nil
}

func validationError(c *fiber.Ctx, err error) error {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var messages []string
		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("%s failed on '%s'", e.Field(), e.Tag()))
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": strings.Join(messages, ", "),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": msg,
	})
}
