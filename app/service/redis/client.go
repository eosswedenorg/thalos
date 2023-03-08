package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewClient(cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb, rdb.Ping(context.Background()).Err()
}
