package redis

import (
	"context"
	"fmt"
	"mpb/configs"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(conf *configs.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Addr,
		DB:   0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{Client: rdb}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
