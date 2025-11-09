package testutils

import (
	"context"
	"mpb/configs"
	"mpb/pkg/db"
	"mpb/pkg/redis"
	"os"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/redis/go-redis/v9"
)

func TestConfig() *configs.Config {
	return &configs.Config{
		Db: configs.DbConfig{
			Dsn: getEnvOrDefault("TEST_DSN", "postgres://mpb:mpb_pas@localhost:5432/mpb_test?sslmode=disable"),
		},
		Redis: configs.RedisConfig{
			Addr: getEnvOrDefault("TEST_REDIS_ADDR", "localhost:6379"),
		},
		JWT: configs.JWTConfig{
			SecretKey:      getEnvOrDefault("TEST_JWT_SECRET", "test-secret-key"),
			AccessTokenTTL: 24 * 60 * 60 * 1000000000,
		},
	}
}

func SetupTestDB(t *testing.T) *db.Db {
	conf := TestConfig()
	database := db.NewDb(conf)

	ctx := context.Background()
	if err := database.Conn.PingContext(ctx); err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return database
}

func SetupTestRedis(t *testing.T) *redis.Redis {
	conf := TestConfig()
	redisClient, err := redis.NewRedis(conf)
	if err != nil {
		t.Fatalf("Failed to connect to test Redis: %v", err)
	}

	ctx := context.Background()
	if err := redisClient.Client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to ping test Redis: %v", err)
	}

	return redisClient
}

func CleanupRedis(t *testing.T, client *redis.Client) {
	ctx := context.Background()
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Logf("Warning: Failed to flush Redis: %v", err)
	}
}

func SetupTestPubSub(t *testing.T) (message.Publisher, message.Subscriber) {
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	return pubsub, pubsub
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
