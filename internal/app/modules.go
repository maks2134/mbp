package app

import (
	"mpb/configs"
	"mpb/internal/auth"
	"mpb/internal/comments"
	"mpb/internal/post_attachments"
	"mpb/internal/posts"
	"mpb/pkg/db"
	"mpb/pkg/redis"
	"mpb/pkg/s3"

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
	// s3 подключить
	s3Client, err := s3.NewS3Client(conf)
	if err != nil {
		logger.Error("failed to init S3 client", err, nil)
		return
	}

	// auth блок
	authRepo := auth.NewAuthRepository(conf, database, redisClient.Client)
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

	// post attachments блоки
	postAttachmentRepo := post_attachments.NewPostAttacmentsRepository(database)
	postAttachmentService := post_attachments.NewPostAttachmentsService(postAttachmentRepo)
	postAttachmentHandler := post_attachments.NewPostAttachmentsHandlers(postAttachmentService, s3Client)
	postAttachmentRoutes := post_attachments.NewPostAttachmentsRoutes(api, postAttachmentHandler, []byte(conf.JWT.SecretKey))
	postAttachmentRoutes.Register()

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
