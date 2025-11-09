package users

import (
	"errors"
	"mpb/internal/posts"
	postsdto "mpb/internal/posts/dto"
	"mpb/internal/users/dto"
	"mpb/pkg/errors_constant"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UsersHandlers struct {
	service *UsersService
}

func NewUsersHandlers(service *UsersService) *UsersHandlers {
	return &UsersHandlers{service: service}
}

// GetUserProfile godoc
// @Summary Get user profile by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 404 {object} fiber.Map
// @Router /api/users/{id} [get]
func (h *UsersHandlers) GetUserProfile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	profile, err := h.service.GetUserProfile(c.Context(), id)
	if err != nil {
		if errors.Is(err, errors_constant.UserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := profileToResponse(profile)
	return c.JSON(response)
}

// GetUserPosts godoc
// @Summary Get posts by user ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} posts.dto.PostResponse
// @Failure 404 {object} fiber.Map
// @Router /api/users/{id}/posts [get]
func (h *UsersHandlers) GetUserPosts(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	filter := posts.PostFilter{
		OnlyActive: true,
	}

	userPosts, err := h.service.GetUserPosts(c.Context(), id, filter)
	if err != nil {
		if errors.Is(err, errors_constant.UserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]postsdto.PostResponse, len(userPosts))
	for i := range userPosts {
		response[i] = convertPostToDTO(&userPosts[i])
	}

	return c.JSON(response)
}

// ListUsers godoc
// @Summary List all users
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserProfileResponse
// @Router /api/users [get]
func (h *UsersHandlers) ListUsers(c *fiber.Ctx) error {
	filter := UserFilter{
		IsActive: func() *bool { b := true; return &b }(),
	}

	users, err := h.service.ListUsers(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]dto.UserProfileResponse, len(users))
	for i, u := range users {
		postsCount, _ := h.service.repo.GetPostsCount(c.Context(), u.ID)
		attachmentsCount, _ := h.service.repo.GetAttachmentsCount(c.Context(), u.ID)

		response[i] = dto.UserProfileResponse{
			ID:               u.ID,
			Name:             u.Name,
			Username:         u.Username,
			Email:            u.Email,
			Age:              u.Age,
			IsActive:         u.IsActive,
			PostsCount:       postsCount,
			AttachmentsCount: attachmentsCount,
			CreatedAt:        u.CreatedAt,
			UpdatedAt:        u.UpdatedAt,
		}
	}

	return c.JSON(response)
}

func profileToResponse(profile *UserProfile) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:               profile.ID,
		Name:             profile.Name,
		Username:         profile.Username,
		Email:            profile.Email,
		Age:              profile.Age,
		IsActive:         profile.IsActive,
		PostsCount:       profile.PostsCount,
		AttachmentsCount: profile.AttachmentsCount,
		CreatedAt:        profile.CreatedAt,
		UpdatedAt:        profile.UpdatedAt,
	}
}

func convertPostToDTO(post *posts.Post) postsdto.PostResponse {
	return postsdto.PostResponse{
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
