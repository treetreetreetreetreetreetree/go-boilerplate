package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-boilerplate/config"

	"github.com/go-redis/redis/v9"
)

const successStatus = "PONG"

var ErrRedisClientConnectionFailed = errors.New("[CACHE] redis client connection was failed")

type Rdb struct {
	Client *redis.Client
}

func Setup(ctx context.Context, cfg config.RedisConfig) (*Rdb, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: "",
		DB:       0,
		PoolSize: 1,
	})

	status := rdb.Ping(ctx)

	if status.Val() != successStatus {
		slog.Error("[CACHE]", "error", ErrRedisClientConnectionFailed)
		return nil, ErrRedisClientConnectionFailed
	}

	slog.Info("[CACHE]", "message", fmt.Sprintf("redis status %s", status))
	return &Rdb{Client: rdb}, nil
}
