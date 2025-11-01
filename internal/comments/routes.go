package comments

import (
	"mpb/internal/comments/dto"
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type CommentsRoutes struct {
	router    fiber.Router
	handler   *CommentsHandlers
	jwtSecret []byte
}

func NewCommentsRoutes(router fiber.Router, handler *CommentsHandlers, jwtSecret []byte) *CommentsRoutes {
	return &CommentsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *CommentsRoutes) Register() {
	comments := r.router.Group("/comments")

	comments.Get("/", r.handler.ListComments)
	comments.Get("/:id", r.handler.GetComment)

	commentsAuth := comments.Group("/", middleware.JWTAuth(r.jwtSecret))

	commentsAuth.Post("/",
		middleware.ValidateBody[dto.CreateCommentRequest](),
		r.handler.CreateComment,
	)

	commentsAuth.Put("/:id",
		middleware.ValidateBody[dto.UpdateCommentRequest](),
		r.handler.UpdateComment,
	)

	commentsAuth.Delete("/:id", r.handler.DeleteComment)
}
