package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}

type DbConfig struct {
	Dsn string
}

type Config struct {
	Db  DbConfig
	JWT JWTConfig
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")

	ttl := 24 * time.Hour
	if v := os.Getenv("JWT_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		JWT: JWTConfig{
			SecretKey:      os.Getenv("JWT_SECRET"),
			AccessTokenTTL: ttl,
		},
	}
}
