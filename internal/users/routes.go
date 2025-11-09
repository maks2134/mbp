package users

import (
	"github.com/gofiber/fiber/v2"
)

type UsersRoutes struct {
	router  fiber.Router
	handler *UsersHandlers
}

func NewUsersRoutes(router fiber.Router, handler *UsersHandlers) *UsersRoutes {
	return &UsersRoutes{
		router:  router,
		handler: handler,
	}
}

func (r *UsersRoutes) Register() {
	users := r.router.Group("/users")

	users.Get("/", r.handler.ListUsers)
	users.Get("/:id", r.handler.GetUserProfile)
	users.Get("/:id/posts", r.handler.GetUserPosts)
}
