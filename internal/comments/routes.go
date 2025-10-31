package comments

import (
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
	//posts := r.router.Group("/posts")

	//posts.Get("/", r.handler.GetAllPosts)

	//auth := posts.Group("/", middleware.JWTAuth(r.jwtSecret))
	//auth.Post("/", r.handler)
	//auth.Put("/:id", middleware.ValidateBody[dto.UpdatePostRequest](), r.handler.UpdatePost)
	//auth.Delete("/:id", r.handler.DeletePost)
}
