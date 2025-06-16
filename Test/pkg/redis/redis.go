package redis

import (
	"context"
	"fmt"
	"myapp/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Conn *redis.Client
}

func New(cfg config.Redis) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - Ping: %w", err)
	}
	return &Redis{rdb}, nil
}

func (r *Redis) CloseConn() {
	r.Conn.Close()
}
