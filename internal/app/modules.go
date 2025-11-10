package app

import (
	"mpb/configs"
	"mpb/internal/auth"
	"mpb/internal/comments"
	"mpb/internal/post_attachments"
	"mpb/internal/posts"
	"mpb/internal/stories"
	"mpb/internal/user_attachments"
	"mpb/internal/users"
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

	// users блок
	usersRepo := users.NewUsersRepository(database)
	usersService := users.NewUsersService(usersRepo, postRepo)
	usersHandler := users.NewUsersHandlers(usersService)
	usersRoutes := users.NewUsersRoutes(api, usersHandler)
	usersRoutes.Register()

	// user attachments блок
	userAttachmentRepo := user_attachments.NewUserAttachmentsRepository(database)
	userAttachmentService := user_attachments.NewUserAttachmentsService(userAttachmentRepo)
	userAttachmentHandler := user_attachments.NewUserAttachmentsHandlers(userAttachmentService, s3Client)
	userAttachmentRoutes := user_attachments.NewUserAttachmentsRoutes(api, userAttachmentHandler, []byte(conf.JWT.SecretKey))
	userAttachmentRoutes.Register()

	// stories блок
	storiesRepo := stories.NewStoriesRepository(database)
	storiesService := stories.NewStoriesService(storiesRepo)
	storiesHandler := stories.NewStoriesHandlers(storiesService, s3Client)
	storiesRoutes := stories.NewStoriesRoutes(api, storiesHandler, []byte(conf.JWT.SecretKey))
	storiesRoutes.Register()

	metricsConsumer := posts.NewMetricsSyncConsumer(postRepo, logger)
	if err := metricsConsumer.StartConsumers(subscriber); err != nil {
		logger.Error("failed to start metrics consumers", err, nil)
	}
}
