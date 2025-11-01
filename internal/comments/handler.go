package comments

import (
	"errors"
	"mpb/internal/comments/dto"
	"mpb/pkg/errors_constant"
	"mpb/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentsHandlers struct {
	service *CommentsService
}

func NewCommentsHandlers(service *CommentsService) *CommentsHandlers {
	return &CommentsHandlers{service: service}
}

// CreateComment godoc
// @Summary Create new comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body dto.CreateCommentRequest true "Comment data"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} map[string]string
// @Router /api/comments [post]
func (h *CommentsHandlers) CreateComment(c *fiber.Ctx) error {
	req := middleware.Body[dto.CreateCommentRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	comment, err := h.service.CreateComment(c.Context(), req.PostID, userID, req.Text)
	if err != nil {
		if errors.Is(err, errors_constant.InvalidCommentText) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toCommentResponse(comment))
}

// GetComment godoc
// @Summary Get comment by ID
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} dto.CommentResponse
// @Failure 404 {object} map[string]string
// @Router /api/comments/{id} [get]
func (h *CommentsHandlers) GetComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment id"})
	}

	comment, err := h.service.repo.FindCommentByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "comment not found"})
	}

	return c.JSON(toCommentResponse(comment))
}

// ListComments godoc
// @Summary List all comments for post
// @Tags Comments
// @Produce json
// @Param post_id query int true "Post ID"
// @Success 200 {array} dto.CommentResponse
// @Router /api/comments [get]
func (h *CommentsHandlers) ListComments(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Query("post_id"))
	if err != nil || postID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post_id"})
	}

	comments, err := h.service.ListComments(c.Context(), postID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp := make([]dto.CommentResponse, len(comments))
	for i, cmt := range comments {
		resp[i] = toCommentResponse(&cmt)
	}
	return c.JSON(resp)
}

// UpdateComment godoc
// @Summary Update comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Updated comment"
// @Success 200 {object} dto.CommentResponse
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/comments/{id} [put]
func (h *CommentsHandlers) UpdateComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment id"})
	}

	req := middleware.Body[dto.UpdateCommentRequest](c)
	if req == nil || req.Text == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	comment, err := h.service.UpdateComment(c.Context(), userID, id, *req.Text)
	if err != nil {
		switch {
		case errors.Is(err, errors_constant.UserNotAuthorized):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can update only your own comments"})
		case errors.Is(err, errors_constant.CommentDeleted):
			return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "comment deleted"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(toCommentResponse(comment))
}

// DeleteComment godoc
// @Summary Delete comment
// @Tags Comments
// @Param id path int true "Comment ID"
// @Success 204 "No Content"
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/comments/{id} [delete]
func (h *CommentsHandlers) DeleteComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	err = h.service.DeleteComment(c.Context(), userID, id)
	if err != nil {
		switch {
		case errors.Is(err, errors_constant.UserNotAuthorized):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can delete only your own comments"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func toCommentResponse(c *Comment) dto.CommentResponse {
	return dto.CommentResponse{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		Text:      c.Text,
		Like:      c.Like,
		Blocked:   c.Blocked,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}
