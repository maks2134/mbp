package stories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"mpb/internal/stories/dto"
	"mpb/pkg/errors_constant"
	"mpb/pkg/s3"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type StoriesHandlers struct {
	service  *StoriesService
	s3Client *s3.S3Client
}

func NewStoriesHandlers(service *StoriesService, s3Client *s3.S3Client) *StoriesHandlers {
	return &StoriesHandlers{service: service, s3Client: s3Client}
}

// CreateStory godoc
// @Summary Create a new story
// @Tags Stories
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Story file (image or video)"
// @Success 201 {object} dto.StoryResponse
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /api/stories [post]
func (h *StoriesHandlers) CreateStory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no file provided"})
	}

	if len(files) > 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "only one file allowed per story"})
	}

	file := files[0]
	ctx := context.Background()

	url, err := h.uploadToS3(ctx, file, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	story, err := h.service.CreateStory(ctx, userID, url, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := storyToResponse(story, false)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *StoriesHandlers) uploadToS3(ctx context.Context, fileHeader *multipart.FileHeader, userID int) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %w", err)
		}
	}(file)

	key := fmt.Sprintf("stories/%d/%d_%s", userID, time.Now().Unix(), fileHeader.Filename)

	url, err := h.s3Client.UploadFile(ctx, key, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}

	return url, nil
}

// GetStory godoc
// @Summary Get story by ID
// @Tags Stories
// @Produce json
// @Param id path int true "Story ID"
// @Success 200 {object} dto.StoryResponse
// @Failure 404 {object} fiber.Map
// @Router /api/stories/{id} [get]
func (h *StoriesHandlers) GetStory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid story id"})
	}

	var viewerUserID *int
	if userID, ok := c.Locals("user_id").(int); ok {
		viewerUserID = &userID
	}

	story, isViewed, err := h.service.GetStory(c.Context(), id, viewerUserID)
	if err != nil {
		if errors.Is(err, errors_constant.UserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "story not found or expired"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := storyToResponse(story, isViewed)
	return c.JSON(response)
}

// ViewStory godoc
// @Summary Mark story as viewed
// @Tags Stories
// @Param id path int true "Story ID"
// @Success 200 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /api/stories/{id}/view [post]
func (h *StoriesHandlers) ViewStory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid story id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if err := h.service.ViewStory(c.Context(), id, userID); err != nil {
		if errors.Is(err, errors_constant.UserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "story not found or expired"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "story viewed"})
}

// ListUserStories godoc
// @Summary Get stories by user ID
// @Tags Stories
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} dto.StoryResponse
// @Router /api/stories/user/{id} [get]
func (h *StoriesHandlers) ListUserStories(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	stories, err := h.service.ListUserStories(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]dto.StoryResponse, len(stories))
	var viewerUserID *int
	if userID, ok := c.Locals("user_id").(int); ok {
		viewerUserID = &userID
	}

	for i, story := range stories {
		var isViewed bool
		if viewerUserID != nil {
			isViewed, _ = h.service.repo.HasUserViewed(c.Context(), story.ID, *viewerUserID)
		}
		response[i] = storyToResponse(&story, isViewed)
	}

	return c.JSON(response)
}

// ListActiveStories godoc
// @Summary Get all active stories
// @Tags Stories
// @Produce json
// @Success 200 {array} dto.StoryResponse
// @Router /api/stories [get]
func (h *StoriesHandlers) ListActiveStories(c *fiber.Ctx) error {
	var excludeUserID *int
	if userID, ok := c.Locals("user_id").(int); ok {
		excludeUserID = &userID
	}

	stories, err := h.service.ListActiveStories(c.Context(), excludeUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]dto.StoryResponse, len(stories))
	var viewerUserID *int
	if userID, ok := c.Locals("user_id").(int); ok {
		viewerUserID = &userID
	}

	for i, story := range stories {
		var isViewed bool
		if viewerUserID != nil {
			isViewed, _ = h.service.repo.HasUserViewed(c.Context(), story.ID, *viewerUserID)
		}
		response[i] = storyToResponse(&story, isViewed)
	}

	return c.JSON(response)
}

// DeleteStory godoc
// @Summary Delete a story
// @Tags Stories
// @Param id path int true "Story ID"
// @Success 204
// @Failure 404 {object} fiber.Map
// @Router /api/stories/{id} [delete]
func (h *StoriesHandlers) DeleteStory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid story id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if err := h.service.DeleteStory(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "story not found or unauthorized"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func storyToResponse(story *Story, isViewed bool) dto.StoryResponse {
	return dto.StoryResponse{
		ID:         story.ID,
		UserID:     story.UserID,
		FileURL:    story.FileURL,
		FileType:   story.FileType,
		ViewsCount: story.ViewsCount,
		ExpiresAt:  story.ExpiresAt,
		CreatedAt:  story.CreatedAt,
		IsViewed:   isViewed,
	}
}
