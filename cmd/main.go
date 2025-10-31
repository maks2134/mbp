package main

import (
	"mpb/configs"
	_ "mpb/docs"
	"mpb/internal/auth"
	"mpb/internal/posts"
	"mpb/pkg/db"
	"mpb/pkg/redis"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title MPB Blog Auth API
// @version 1.0
// @description Authentication service for MPB blog platform
// @host localhost:8000
// @BasePath /api
func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)

	redisClient, err := redis.NewRedis(conf)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	app := fiber.New()

	api := app.Group("/api")

	//регаем watermill, для синхронизации redis и postgres
	logger := watermill.NewStdLogger(false, false)
	publisher := gochannel.NewGoChannel(gochannel.Config{}, logger)
	subscriber := gochannel.NewGoChannel(gochannel.Config{}, logger)

	// auth блок
	authRepo := auth.NewAuthRepository(conf, database)
	authService := auth.NewAuthService(authRepo, []byte(conf.JWT.SecretKey), conf.JWT.AccessTokenTTL)
	authHandler := auth.NewAuthHandlers(authService)
	authRoutes := auth.NewAuthRoutes(api, authHandler)
	authRoutes.Register()

	// posts блок
	postRepo := posts.NewPostsRepository(database)
	metricsService := posts.NewMetricsService(redisClient.Client, publisher, logger)
	postService := posts.NewPostsService(postRepo, metricsService, publisher, logger)
	postsHandler := posts.NewPostsHandlers(postService, metricsService)
	postsRoutes := posts.NewPostsRoutes(api, postsHandler, []byte(conf.JWT.SecretKey))
	postsRoutes.Register()

	metricsConsumer := posts.NewMetricsSyncConsumer(postRepo, logger)
	if err := metricsConsumer.StartConsumers(subscriber); err != nil {
		log.Fatalf("Failed to start metrics consumers: %v", err)
	}

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
