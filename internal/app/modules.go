package app

import (
	"mpb/configs"
	"mpb/internal/auth"
	"mpb/internal/comments"
	"mpb/internal/posts"
	"mpb/pkg/db"
	"mpb/pkg/redis"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2"
)

func RegisterModules(
	api fiber.Router,
	database *db.Db,
	redisClient *redis.Redis,
	publisher message.Publisher,
	subscriber message.Subscriber,
	logger watermill.LoggerAdapter,
	conf *configs.Config,
) {
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
	postHandler := posts.NewPostsHandlers(postService, metricsService)
	postRoutes := posts.NewPostsRoutes(api, postHandler, []byte(conf.JWT.SecretKey))
	postRoutes.Register()

	// comments блок
	commentRepo := comments.NewCommentsRepository(database)
	commentService := comments.NewCommentsService(commentRepo)
	commentHandler := comments.NewCommentsHandlers(commentService)
	commentRoutes := comments.NewCommentsRoutes(api, commentHandler, []byte(conf.JWT.SecretKey))
	commentRoutes.Register()

	metricsConsumer := posts.NewMetricsSyncConsumer(postRepo, logger)
	if err := metricsConsumer.StartConsumers(subscriber); err != nil {
		logger.Error("failed to start metrics consumers", err, nil)
	}
}
