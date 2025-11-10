package user_attachments

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"mpb/pkg/errors_constant"
	"mpb/pkg/s3"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserAttachmentsHandlers struct {
	service  *UserAttachmentsService
	s3Client *s3.S3Client
}

func NewUserAttachmentsHandlers(service *UserAttachmentsService, s3Client *s3.S3Client) *UserAttachmentsHandlers {
	return &UserAttachmentsHandlers{service: service, s3Client: s3Client}
}

// UploadAttachments godoc
// @Summary Upload attachments for a user profile
// @Tags UserAttachments
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param files formData file true "Files to upload" collectionFormat(multi)
// @Success 201 {array} UserAttachment
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id}/attachments [post]
func (h *UserAttachmentsHandlers) UploadAttachments(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil || userID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	// Check if user is authorized (can only upload to own profile)
	authenticatedUserID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if authenticatedUserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can only upload attachments to your own profile"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no files provided"})
	}

	var uploaded []UserAttachment
	ctx := context.Background()

	for _, file := range files {
		url, err := h.uploadToS3(ctx, file, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		attachment := &UserAttachment{
			UserID:   userID,
			FileURL:  url,
			FileType: file.Header.Get("Content-Type"),
			FileSize: file.Size,
		}

		if err := h.service.CreateAttachment(ctx, attachment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		uploaded = append(uploaded, *attachment)
	}

	return c.Status(fiber.StatusCreated).JSON(uploaded)
}

func (h *UserAttachmentsHandlers) uploadToS3(ctx context.Context, fileHeader *multipart.FileHeader, userID int) (string, error) {
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

	key := fmt.Sprintf("users/%d/%s", userID, fileHeader.Filename)

	url, err := h.s3Client.UploadFile(ctx, key, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}

	return url, nil
}

// GetAttachments godoc
// @Summary Get attachments for a user
// @Tags UserAttachments
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} UserAttachment
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id}/attachments [get]
func (h *UserAttachmentsHandlers) GetAttachments(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	attachments, err := h.service.ListAttachments(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(attachments)
}

// DeleteAttachment godoc
// @Summary Delete a user attachment
// @Tags UserAttachments
// @Param id path int true "Attachment ID"
// @Success 204
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/attachments/{id} [delete]
func (h *UserAttachmentsHandlers) DeleteAttachment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid attachment id"})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	if err := h.service.DeleteAttachment(c.Context(), id, userID); err != nil {
		if errors.Is(err, errors_constant.UserNotAuthorized) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can only delete your own attachments"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
