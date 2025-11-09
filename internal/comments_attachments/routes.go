package comments_attachments

import (
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type CommentAttachmentsRoutes struct {
	router    fiber.Router
	handler   *CommentAttachmentsHandlers
	jwtSecret []byte
}

func NewCommentAttachmentsRoutes(router fiber.Router, handler *CommentAttachmentsHandlers, jwtSecret []byte) *CommentAttachmentsRoutes {
	return &CommentAttachmentsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *CommentAttachmentsRoutes) Register() {
	group := r.router.Group("/comments")

	group.Get("/:id/attachments", r.handler.GetAttachments)

	authGroup := group.Group("/", middleware.JWTAuth(r.jwtSecret))
	authGroup.Post("/:id/attachments", r.handler.UploadAttachments)
	authGroup.Delete("/attachments/:id", r.handler.DeleteAttachment)
}
