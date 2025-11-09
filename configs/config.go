package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AWSConfig struct {
	Region string
	Bucket string
}

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}

type DbConfig struct {
	Dsn string
}

type RedisConfig struct {
	Addr string
}

type Config struct {
	Db    DbConfig
	Redis RedisConfig
	JWT   JWTConfig
	AWS   AWSConfig
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")

	ttl := 24 * time.Hour
	if v := os.Getenv("JWT_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Redis: RedisConfig{
			Addr: redisAddr,
		},
		JWT: JWTConfig{
			SecretKey:      os.Getenv("JWT_SECRET"),
			AccessTokenTTL: ttl,
		},
		AWS: AWSConfig{
			Region: os.Getenv("AWS_REGION"),
			Bucket: os.Getenv("AWS_BUCKET"),
		},
	}
}
