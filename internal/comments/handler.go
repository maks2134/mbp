package comments

import (
	"github.com/gofiber/fiber/v2"
)

type CommentsHandler struct {
	// CommentsHandler *CommentsHandler
}

func NewCommentsHandler() *CommentsHandler {
	return &CommentsHandler{}
}

func (h *CommentsHandler) GetAllComments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *CommentsHandler) CreateComments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *CommentsHandler) DeleteComments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *CommentsHandler) UpdateComments(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
