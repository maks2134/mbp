package posts

import "github.com/gofiber/fiber/v2"

type PostsHandlers struct {
	postsService *PostsService
}

func NewPostsHandlers(service *PostsService) *PostsHandlers {
	return &PostsHandlers{
		postsService: service,
	}
}

func (h *PostsHandlers) SavePosts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *PostsHandlers) GetPosts(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *PostsHandlers) GetPost(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *PostsHandlers) DeletePost(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *PostsHandlers) UpdatePost(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
