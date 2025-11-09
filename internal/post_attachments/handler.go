package post_attachments

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"mpb/pkg/s3"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostAttachmentsHandlers struct {
	service  *PostAttachmentsService
	s3Client *s3.S3Client
}

func NewPostAttachmentsHandlers(service *PostAttachmentsService, s3Client *s3.S3Client) *PostAttachmentsHandlers {
	return &PostAttachmentsHandlers{service: service, s3Client: s3Client}
}

// UploadAttachments godoc
// @Summary Upload attachments for a post
// @Tags PostAttachments
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Post ID"
// @Param files formData file true "Files to upload" collectionFormat(multi)
// @Success 201 {array} PostAttachment
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /posts/{id}/attachments [post]
func (h *PostAttachmentsHandlers) UploadAttachments(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil || postID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no files provided"})
	}

	var uploaded []PostAttachment
	ctx := context.Background()

	for _, file := range files {
		url, err := h.uploadToS3(ctx, file, postID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		attachment := &PostAttachment{
			PostID:   postID,
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

func (h *PostAttachmentsHandlers) uploadToS3(ctx context.Context, fileHeader *multipart.FileHeader, postID int) (string, error) {
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

	key := fmt.Sprintf("posts/%d/%s", postID, fileHeader.Filename)

	url, err := h.s3Client.UploadFile(ctx, key, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}

	return url, nil
}

// GetAttachments godoc
// @Summary Get attachments for a post
// @Tags PostAttachments
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} PostAttachment
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /posts/{id}/attachments [get]
func (h *PostAttachmentsHandlers) GetAttachments(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post id"})
	}

	attachments, err := h.service.ListAttachments(c.Context(), postID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(attachments)
}

// DeleteAttachment godoc
// @Summary Delete a post attachment
// @Tags PostAttachments
// @Param id path int true "Attachment ID"
// @Success 204
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /posts/attachments/{id} [delete]
func (h *PostAttachmentsHandlers) DeleteAttachment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid attachment id"})
	}

	if err := h.service.DeleteAttachment(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
