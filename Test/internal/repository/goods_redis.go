package repository

import (
	"context"
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"myapp/pkg/redis"
	"time"
)

func NewGoodsRedis(rdb *redis.Redis) *GoodsRedis {
	return &GoodsRedis{
		Rdb: rdb,
	}
}

type GoodsRedis struct {
	Rdb *redis.Redis
}

func (db *GoodsRedis) GetRoleRights(ctx context.Context, role string) ([]string, error) {
	rights, err := db.Rdb.Conn.SMembers(ctx, role).Result()
	if err != nil {
		if err == redis2.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("AuthRedis - GetRoleRights - db.Rdb.SMembers: %w", err)
	}

	return rights, nil
}

func (db *GoodsRedis) AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error {
	err := db.Rdb.Conn.SAdd(ctx, role, rights).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.SAdd: %w", err)
	}
	err = db.Rdb.Conn.Expire(ctx, role, period).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.Expire: %w", err)
	}

	return nil
}
