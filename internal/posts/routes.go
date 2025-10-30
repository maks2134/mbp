package posts

import (
	"mpb/internal/posts/dto"
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type PostsRoutes struct {
	router    fiber.Router
	handler   *PostsHandlers
	jwtSecret []byte
}

func NewPostsRoutes(router fiber.Router, handler *PostsHandlers, jwtSecret []byte) *PostsRoutes {
	return &PostsRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *PostsRoutes) Register() {
	posts := r.router.Group("/posts")

	posts.Get("/", r.handler.GetAllPosts)
	posts.Get("/:id", r.handler.GetPost)

	auth := posts.Group("/", middleware.JWTAuth(r.jwtSecret))
	auth.Post("/", middleware.ValidateBody[dto.CreatePostRequest](), r.handler.CreatePost)
	auth.Put("/:id", middleware.ValidateBody[dto.UpdatePostRequest](), r.handler.UpdatePost)
	auth.Delete("/:id", r.handler.DeletePost)
}
