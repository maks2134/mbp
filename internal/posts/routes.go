package posts

import (
	"github.com/gofiber/fiber/v2"
)

type PostsRoutes struct {
	router  fiber.Router
	handler *PostsHandlers
}

func NewPostsRoutes(router fiber.Router, handlers *PostsHandlers) *PostsRoutes {
	return &PostsRoutes{
		router:  router,
		handler: handlers,
	}
}

func (r *PostsRoutes) Register() {
	auth := r.router.Group("/posts")

	auth.Post("/", r.handler.SavePosts)
	auth.Get("/", r.handler.GetPosts)
	auth.Get("/:id", r.handler.GetPost)
	auth.Put("/:id", r.handler.UpdatePost)
	auth.Delete("/:id", r.handler.DeletePost)
}
