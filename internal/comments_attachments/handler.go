package comments_attachments

import (
	"context"
	"mime/multipart"
	"mpb/pkg/s3"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentAttachmentsHandlers struct {
	service  *CommentAttachmentsService
	s3Client *s3.S3Client
}

func NewCommentAttachmentsHandlers(service *CommentAttachmentsService, s3Client *s3.S3Client) *CommentAttachmentsHandlers {
	return &CommentAttachmentsHandlers{service: service, s3Client: s3Client}
}

func (h *CommentAttachmentsHandlers) UploadAttachments(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no files provided"})
	}

	var uploaded []CommentAttachment
	ctx := context.Background()

	for _, file := range files {
		url, err := h.uploadToS3(ctx, file, commentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		attachment := &CommentAttachment{
			CommentID: commentID,
			FileURL:   url,
			FileType:  file.Header.Get("Content-Type"),
			FileSize:  file.Size,
		}

		if err := h.service.CreateAttachment(ctx, attachment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		uploaded = append(uploaded, *attachment)
	}

	return c.Status(fiber.StatusCreated).JSON(uploaded)
}

func (h *CommentAttachmentsHandlers) uploadToS3(ctx context.Context, file *multipart.FileHeader, commentID int) (string, error) {
	path := "comments/" + strconv.Itoa(commentID)
	return h.s3Client.UploadFile(ctx, file, path)
}

func (h *CommentAttachmentsHandlers) GetAttachments(c *fiber.Ctx) error {
	commentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid comment id"})
	}

	attachments, err := h.service.ListAttachments(c.Context(), commentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(attachments)
}

func (h *CommentAttachmentsHandlers) DeleteAttachment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid attachment id"})
	}

	if err := h.service.DeleteAttachment(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
