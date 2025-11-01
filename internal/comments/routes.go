package comments

import (
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type CommentsRoutes struct {
	router    fiber.Router
	handler   *CommentsHandler
	jwtSecret []byte
}

func NewCommentsRoutes(router fiber.Router, handler *CommentsHandler, jwtSecret []byte) *CommentsRoutes {
	return &CommentsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *CommentsRoutes) Register() {
	comments := r.router.Group("/comments")

	comments.Get("/", r.handler.GetAllComments)

	commentsAuth := comments.Group("/", middleware.JWTAuth(r.jwtSecret))
	commentsAuth.Post("/", r.handler.CreateComments)
	commentsAuth.Put("/:id", r.handler.UpdateComments)
	commentsAuth.Delete("/:id", r.handler.DeleteComments)
}
