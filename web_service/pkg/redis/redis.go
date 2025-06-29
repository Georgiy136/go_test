package redis

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/web_service/config"

	"github.com/redis/go-redis/v9"
)

func NewConn(cfg config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - Ping: %w", err)
	}
	return rdb, nil
}
