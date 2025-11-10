package stories

import (
	"mpb/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type StoriesRoutes struct {
	router    fiber.Router
	handler   *StoriesHandlers
	jwtSecret []byte
}

func NewStoriesRoutes(router fiber.Router, handler *StoriesHandlers, jwtSecret []byte) *StoriesRoutes {
	return &StoriesRoutes{router: router, handler: handler, jwtSecret: jwtSecret}
}

func (r *StoriesRoutes) Register() {
	stories := r.router.Group("/stories")

	stories.Get("/", r.handler.ListActiveStories)
	stories.Get("/:id", r.handler.GetStory)
	stories.Get("/user/:id", r.handler.ListUserStories)

	authGroup := stories.Group("/", middleware.JWTAuth(r.jwtSecret))
	authGroup.Post("/", r.handler.CreateStory)
	authGroup.Post("/:id/view", r.handler.ViewStory)
	authGroup.Delete("/:id", r.handler.DeleteStory)
}
