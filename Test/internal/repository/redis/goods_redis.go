package redis

import (
	"context"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"myapp/internal/models"
	"myapp/internal/usecase"
	connect "myapp/pkg/redis"
	"time"
)

func NewGoodsRedis(rdb *connect.Redis) usecase.GoodsCache {
	return &GoodsRedis{rdb}
}

type GoodsRedis struct {
	*connect.Redis
}

const (
	goodsKeyFormat = "%d_%d"
	expPeriod      = 1 * time.Minute
)

func (db *GoodsRedis) GetGoods(ctx context.Context, goodsID, projectID int) (*models.Goods, error) {
	if db.Conn == nil {
		return nil, nil
	}

	key := fmt.Sprintf(goodsKeyFormat, goodsID, projectID)

	bytesString, err := db.Conn.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("GetGoods - db.Conn.Get: %w", err)
	}
	goods := models.Goods{}
	if err = jsoniter.UnmarshalFromString(bytesString, goods); err != nil {
		return nil, fmt.Errorf("GetGoods - jsoniter.UnmarshalFromString err: %w", err)
	}
	return &goods, nil
}

func (db *GoodsRedis) SaveGoods(ctx context.Context, goodsID, projectID int, goods models.Goods) error {
	if db.Conn == nil {
		return nil
	}

	key := fmt.Sprintf(goodsKeyFormat, goodsID, projectID)

	bytesData, err := json.Marshal(goods)
	if err != nil {
		return fmt.Errorf("SaveGoods - json.Marshal err: %w", err)
	}

	err = db.Conn.Append(ctx, key, string(bytesData)).Err()
	if err != nil {
		return fmt.Errorf("SaveGoods - db.Conn.Append: %w", err)
	}
	err = db.Conn.Expire(ctx, key, expPeriod).Err()
	if err != nil {
		return fmt.Errorf("SaveGoods - db.Rdb.Expire: %w", err)
	}

	return nil
}
