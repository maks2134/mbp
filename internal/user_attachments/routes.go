package user_attachments

import (
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserAttachmentsRoutes struct {
	router    fiber.Router
	handler   *UserAttachmentsHandlers
	jwtSecret []byte
}

func NewUserAttachmentsRoutes(router fiber.Router, handler *UserAttachmentsHandlers, jwtSecret []byte) *UserAttachmentsRoutes {
	return &UserAttachmentsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *UserAttachmentsRoutes) Register() {
	group := r.router.Group("/users")

	group.Get("/:id/attachments", r.handler.GetAttachments)

	authGroup := group.Group("/", middleware.JWTAuth(r.jwtSecret))
	authGroup.Post("/:id/attachments", r.handler.UploadAttachments)
	authGroup.Delete("/attachments/:id", r.handler.DeleteAttachment)
}
