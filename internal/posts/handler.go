package posts

import (
	"errors"
	"mpb/internal/posts/dto"
	"mpb/pkg/errors_constant"
	"mpb/pkg/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostsHandlers struct {
	service *PostsService
}

func NewPostsHandlers(service *PostsService) *PostsHandlers {
	return &PostsHandlers{service: service}
}

// CreatePost godoc
// @Summary Create new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post data"
// @Success 201 {object} dto.PostResponse
// @Failure 400 {object} map[string]string
// @Router /api/posts [post]
func (h *PostsHandlers) CreatePost(c *fiber.Ctx) error {
	req := middleware.Body[dto.CreatePostRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID := c.Locals("user_id").(int)

	post, err := h.service.CreatePost(c.Context(), userID, *req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetPost godoc
// @Summary Get post by ID
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} dto.PostResponse
// @Failure 404 {object} map[string]string
// @Router /api/posts/{id} [get]
func (h *PostsHandlers) GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	post, err := h.service.GetPostByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, errors_constant.NotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(post)
}

// GetAllPosts godoc
// @Summary List all posts
// @Tags Posts
// @Produce json
// @Success 200 {array} dto.PostResponse
// @Router /api/posts [get]
func (h *PostsHandlers) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.service.GetAllPosts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New("error getting all posts").Error()})
	}
	return c.JSON(posts)
}

// UpdatePost godoc
// @Summary Update existing post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body dto.UpdatePostRequest true "Updated post data"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} map[string]string
// @Router /api/posts/{id} [put]
func (h *PostsHandlers) UpdatePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	req := middleware.Body[dto.UpdatePostRequest](c)
	if req == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	userID := c.Locals("user_id").(int)

	post, err := h.service.UpdatePost(c.Context(), userID, id, *req)
	if err != nil {
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(post)
}

// DeletePost godoc
// @Summary Delete post
// @Tags Posts
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /api/posts/{id} [delete]
func (h *PostsHandlers) DeletePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	userID := c.Locals("user_id").(int)

	err = h.service.DeletePost(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
