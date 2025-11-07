package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	SecretKey             string
	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	RefreshTokenSecretKey string
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
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")

	accessTtl := 24 * time.Hour
	if v := os.Getenv("JWT_ACCESS_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			accessTtl = d
		}
	}

	refreshTtl := 720 * time.Hour // 30 дней
	if v := os.Getenv("JWT_REFRESH_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			refreshTtl = d
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
			SecretKey:             os.Getenv("JWT_SECRET"),
			AccessTokenTTL:        accessTtl,
			RefreshTokenTTL:       refreshTtl,
			RefreshTokenSecretKey: os.Getenv("JWT_REFRESH_SECRET"),
		},
	}
}
