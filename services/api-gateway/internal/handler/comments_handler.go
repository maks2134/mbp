package handler

import (
	posts_proto "mpb/proto/posts"
	"mpb/services/api-gateway/internal/client"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentsHandler struct {
	client    *client.CommentsClient
	jwtSecret []byte
}

func NewCommentsHandler(client *client.CommentsClient, jwtSecret []byte) *CommentsHandler {
	return &CommentsHandler{
		client:    client,
		jwtSecret: jwtSecret,
	}
}

func (h *CommentsHandler) CreateComment(c *fiber.Ctx) error {
	// TODO: Add JWT authentication middleware
	userID := 1
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid post id"})
	}

	var req posts_proto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	req.PostId = int32(postID)
	req.UserId = int32(userID)

	resp, err := h.client.CreateComment(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(resp.Comment)
}

func (h *CommentsHandler) ListComments(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid post id"})
	}

	req := &posts_proto.ListCommentsRequest{PostId: int32(postID)}
	resp, err := h.client.ListComments(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp.Comments)
}

func (h *CommentsHandler) GetComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid comment id"})
	}

	req := &posts_proto.GetCommentRequest{CommentId: int32(id)}
	resp, err := h.client.GetComment(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp.Comment)
}

func (h *CommentsHandler) UpdateComment(c *fiber.Ctx) error {
	// TODO: Implement
	return c.Status(501).JSON(fiber.Map{"error": "not implemented"})
}

func (h *CommentsHandler) DeleteComment(c *fiber.Ctx) error {
	// TODO: Implement
	return c.Status(501).JSON(fiber.Map{"error": "not implemented"})
}
