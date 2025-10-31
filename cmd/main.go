package main

import (
	"mpb/configs"
	_ "mpb/docs"
	"mpb/internal/auth"
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
	defer redisClient.Close()

	app := fiber.New()

	api := app.Group("/api")

	// Настройка Watermill для event-driven архитектуры
	// КРИТИЧНО: В gochannel publisher и subscriber ДОЛЖНЫ быть одним и тем же экземпляром
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	// Используем один и тот же экземпляр как publisher и subscriber
	// Это критично для работы gochannel - разные экземпляры не видят подписчиков друг друга
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

	// Запускаем consumer для синхронизации метрик из Redis в PostgreSQL
	// КРИТИЧНО: запускаем ПЕРЕД стартом сервера и ждем полной инициализации
	metricsConsumer := posts.NewMetricsSyncConsumer(postRepo, logger)
	if err := metricsConsumer.StartConsumers(subscriber); err != nil {
		log.Fatalf("Failed to start metrics consumers: %v", err)
	}

	// Даем дополнительное время для полной регистрации подписок в gochannel
	// Принудительно переключаем контекст чтобы горутины точно запустились и достигли range
	for i := 0; i < 20; i++ {
		runtime.Gosched() // Даем горутинам время на запуск и достижение range
	}
	time.Sleep(500 * time.Millisecond)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Info("Starting server on :8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
