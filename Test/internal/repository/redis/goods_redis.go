package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	connect "myapp/pkg/redis"
	"time"
)

func NewGoodsRedis(rdb *connect.Redis) *GoodsRedis {
	return &GoodsRedis{rdb}
}

type GoodsRedis struct {
	*connect.Redis
}

func (db *GoodsRedis) GetRoleRights(ctx context.Context, role string) ([]string, error) {
	rights, err := db.Conn.SMembers(ctx, role).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("AuthRedis - GetRoleRights - db.Rdb.SMembers: %w", err)
	}

	return rights, nil
}

func (db *GoodsRedis) AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error {
	err := db.Conn.SAdd(ctx, role, rights).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.SAdd: %w", err)
	}
	err = db.Conn.Expire(ctx, role, period).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.Expire: %w", err)
	}

	return nil
}
