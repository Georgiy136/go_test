package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewGoodsRedis(rdb *redis.Client) *GoodsRedis {
	return &GoodsRedis{
		Rdb: rdb,
	}
}

type GoodsRedis struct {
	Rdb *redis.Client
}

func (db *GoodsRedis) GetRoleRights(ctx context.Context, role string) ([]string, error) {
	rights, err := db.Rdb.SMembers(ctx, role).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("AuthRedis - GetRoleRights - db.Rdb.SMembers: %w", err)
	}

	return rights, nil
}

func (db *GoodsRedis) AddRoleRights(ctx context.Context, role string, rights []string, period time.Duration) error {
	err := db.Rdb.SAdd(ctx, role, rights).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.SAdd: %w", err)
	}
	err = db.Rdb.Expire(ctx, role, period).Err()
	if err != nil {
		return fmt.Errorf("AuthRedis - AddRoleRights - db.Rdb.Expire: %w", err)
	}

	return nil
}
