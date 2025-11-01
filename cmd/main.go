package main

import (
	"mpb/configs"
	_ "mpb/docs"
	"mpb/internal/auth"
	"mpb/internal/comments"
	"mpb/internal/posts"
	"mpb/pkg/db"
	"mpb/pkg/redis"
	"runtime"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
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
	defer func(redisClient *redis.Redis) {
		err := redisClient.Close()
		if err != nil {
			log.Fatalf("Failed to close Redis connection: %v", err)
		}
	}(redisClient)

	app := fiber.New()

	api := app.Group("/api")

	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	publisher := pubsub
	subscriber := pubsub

	log.Infof("Watermill pubsub initialized, publisher=%p, subscriber=%p",
		message.Publisher(publisher), message.Subscriber(subscriber))

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

	// comments блок
	commentRepo := comments.NewCommentsRepository(database)
	commentsService := comments.NewCommentsService(commentRepo)
	commentsHandler := comments.NewCommentsHandlers(commentsService)
	commentsRoutes := comments.NewCommentsRoutes(api, commentsHandler, []byte(conf.JWT.SecretKey))
	commentsRoutes.Register()

	metricsConsumer := posts.NewMetricsSyncConsumer(postRepo, logger)
	if err := metricsConsumer.StartConsumers(subscriber); err != nil {
		log.Fatalf("Failed to start metrics consumers: %v", err)
	}

	for i := 0; i < 20; i++ {
		runtime.Gosched()
	}
	time.Sleep(500 * time.Millisecond)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Info("Starting server on :8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
