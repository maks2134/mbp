package post_attachments

import (
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type PostAttachmentsRoutes struct {
	router    fiber.Router
	handler   *PostAttachmentsHandlers
	jwtSecret []byte
}

func NewPostAttachmentsRoutes(router fiber.Router, handler *PostAttachmentsHandlers, jwtSecret []byte) *PostAttachmentsRoutes {
	return &PostAttachmentsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *PostAttachmentsRoutes) Register() {
	group := r.router.Group("/posts")

	group.Get("/:id/attachments", r.handler.GetAttachments)

	authGroup := group.Group("/", middleware.JWTAuth(r.jwtSecret))
	authGroup.Post("/:id/attachments", r.handler.UploadAttachments)
	authGroup.Delete("/attachments/:id", r.handler.DeleteAttachment)
}
