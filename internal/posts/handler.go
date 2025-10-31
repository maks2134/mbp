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
	service        *PostsService
	metricsService *MetricsService
}

func NewPostsHandlers(service *PostsService, metricsService *MetricsService) *PostsHandlers {
	return &PostsHandlers{
		service:        service,
		metricsService: metricsService,
	}
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

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	post, err := h.service.CreatePost(c.Context(), userID, req.Title, req.Description, req.Tag)
	if err != nil {
		if errors.Is(err, errors_constant.InvalidTitle) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := postToResponse(post)
	return c.Status(fiber.StatusCreated).JSON(response)
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
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := postToResponse(post)
	return c.JSON(response)
}

// GetAllPosts godoc
// @Summary List all posts
// @Tags Posts
// @Produce json
// @Success 200 {array} dto.PostResponse
// @Router /api/posts [get]
func (h *PostsHandlers) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.service.ListPosts(c.Context(), PostFilter{OnlyActive: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		response[i] = postToResponse(&post)
	}
	return c.JSON(response)
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

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	currentPost, err := h.service.GetPostByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if currentPost.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can update only your own posts"})
	}

	title := currentPost.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := currentPost.Description
	if req.Description != nil {
		description = *req.Description
	}
	tag := currentPost.Tag
	if req.Tag != nil {
		tag = *req.Tag
	}

	post, err := h.service.UpdatePost(c.Context(), userID, id, title, description, tag)
	if err != nil {
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		if errors.Is(err, errors_constant.UserNotAuthorized) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can update only your own posts"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := postToResponse(post)
	return c.JSON(response)
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

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	err = h.service.DeletePost(c.Context(), userID, id)
	if err != nil {
		if errors.Is(err, errors_constant.PostNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		if errors.Is(err, errors_constant.UserNotAuthorized) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can delete only your own posts"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// LikePost godoc
// @Summary Like a post
// @Tags Posts
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/posts/{id}/like [post]
func (h *PostsHandlers) LikePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if err := h.metricsService.LikePost(c.Context(), userID, id); err != nil {
		if err.Error() == "user already liked this post" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	likes, err := h.metricsService.GetLikes(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"post_id": id,
		"likes":   likes,
		"message": "Post liked successfully",
	})
}

// UnlikePost godoc
// @Summary Unlike a post
// @Tags Posts
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/posts/{id}/unlike [delete]
func (h *PostsHandlers) UnlikePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if err := h.metricsService.UnlikePost(c.Context(), userID, id); err != nil {
		if err.Error() == "user hasn't liked this post" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	likes, err := h.metricsService.GetLikes(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"post_id": id,
		"likes":   likes,
		"message": "Post unliked successfully",
	})
}

func postToResponse(post *Post) dto.PostResponse {
	return dto.PostResponse{
		ID:           post.ID,
		UserID:       post.UserID,
		Title:        post.Title,
		Description:  post.Description,
		Tag:          post.Tag,
		Like:         post.Like,
		CountViewers: post.CountViewers,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}
}
