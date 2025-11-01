package app

import (
	"context"
	"fmt"
	"log"
	"mpb/configs"
	"mpb/pkg/db"
	"mpb/pkg/redis"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type App struct {
	conf       *configs.Config
	db         *db.Db
	redis      *redis.Redis
	fiberApp   *fiber.App
	publisher  message.Publisher
	subscriber message.Subscriber
	logger     watermill.LoggerAdapter
}

func New(conf *configs.Config) (*App, error) {
	database := db.NewDb(conf)

	redisClient, err := redis.NewRedis(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	app := fiber.New()

	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	return &App{
		conf:       conf,
		db:         database,
		redis:      redisClient,
		fiberApp:   app,
		publisher:  pubsub,
		subscriber: pubsub,
		logger:     logger,
	}, nil
}

func (a *App) Run() error {
	api := a.fiberApp.Group("/api")
	RegisterModules(api, a.db, a.redis, a.publisher, a.subscriber, a.logger, a.conf)

	a.fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	for i := 0; i < 20; i++ {
		runtime.Gosched()
	}
	time.Sleep(500 * time.Millisecond)

	port := ":8000"
	log.Printf("Starting server on %s", port)

	serverErrCh := make(chan error, 1)
	go func() {
		if err := a.fiberApp.Listen(port); err != nil {
			serverErrCh <- err
		} else {
			serverErrCh <- nil
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("Received signal %s. Shutting down...", sig.String())
	case err := <-serverErrCh:
		if err != nil {
			log.Printf("HTTP server stopped with error: %v", err)
		} else {
			log.Printf("HTTP server stopped.")
		}
	}

	shutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.fiberApp.Shutdown(); err != nil {
		log.Printf("Error during fiber shutdown: %v", err)
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			log.Printf("Error closing redis: %v", err)
		}
	}

	if a.db != nil && a.db.Conn.Close() != nil {
		if err := a.db.Conn.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	<-ctx.Done()
	log.Printf("Shutdown complete")
	return nil
}
