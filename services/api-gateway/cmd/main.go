package main

import (
	"log"
	"mpb/configs"
	"mpb/services/api-gateway/internal/client"
	"mpb/services/api-gateway/internal/router"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	conf := configs.LoadConfig()

	// Initialize gRPC clients
	postsServiceAddr := os.Getenv("POSTS_SERVICE_ADDR")
	if postsServiceAddr == "" {
		postsServiceAddr = "localhost:50051"
	}

	storiesServiceAddr := os.Getenv("STORIES_SERVICE_ADDR")
	if storiesServiceAddr == "" {
		storiesServiceAddr = "localhost:50052"
	}

	usersServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if usersServiceAddr == "" {
		usersServiceAddr = "localhost:50053"
	}

	postsClient, err := client.NewPostsClient(postsServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create posts client: %v", err)
	}
	defer postsClient.Close()

	commentsClient, err := client.NewCommentsClient(postsServiceAddr) // Comments service is in posts-service
	if err != nil {
		log.Fatalf("Failed to create comments client: %v", err)
	}
	defer commentsClient.Close()

	storiesClient, err := client.NewStoriesClient(storiesServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create stories client: %v", err)
	}
	defer storiesClient.Close()

	usersClient, err := client.NewUsersClient(usersServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create users client: %v", err)
	}
	defer usersClient.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	api := app.Group("/api")
	router.SetupRoutes(api, postsClient, commentsClient, storiesClient, usersClient, []byte(conf.JWT.SecretKey))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}
		log.Printf("API Gateway starting on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-sigCh
	log.Println("Shutting down API Gateway...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("API Gateway stopped")
}
