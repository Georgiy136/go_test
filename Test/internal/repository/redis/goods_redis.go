package redis

import (
	"context"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"time"
)

func NewGoodsRedis(conn *redis.Client) usecase.GoodsCache {
	return &GoodsRedis{conn: conn}
}

type GoodsRedis struct {
	conn *redis.Client
}

const (
	goodsKeyFormat = "%d_%d"
	expPeriod      = 1 * time.Minute
)

func (cache *GoodsRedis) GetGoods(ctx context.Context, goodsID, projectID int) (*models.Goods, error) {
	if cache.conn != nil {
		key := fmt.Sprintf(goodsKeyFormat, goodsID, projectID)

		bytesString, err := cache.conn.Get(ctx, key).Result()
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
	return nil, nil
}

func (cache *GoodsRedis) SaveGoods(ctx context.Context, goodsID, projectID int, goods models.Goods) error {
	if cache.conn != nil {
		key := fmt.Sprintf(goodsKeyFormat, goodsID, projectID)

		bytesData, err := json.Marshal(goods)
		if err != nil {
			return fmt.Errorf("SaveGoods - json.Marshal err: %w", err)
		}

		err = cache.conn.Append(ctx, key, string(bytesData)).Err()
		if err != nil {
			return fmt.Errorf("SaveGoods - db.Conn.Append: %w", err)
		}
		err = cache.conn.Expire(ctx, key, expPeriod).Err()
		if err != nil {
			return fmt.Errorf("SaveGoods - db.Rdb.Expire: %w", err)
		}
	}
	return nil
}

func (cache *GoodsRedis) ClearGoods(ctx context.Context, goodsID, projectID int) error {
	if cache.conn != nil {
		key := fmt.Sprintf(goodsKeyFormat, goodsID, projectID)
		_, err := cache.conn.Del(ctx, key).Result()
		if err != nil {
			return fmt.Errorf("GetGoods - db.Conn.Get: %w", err)
		}
	}
	return nil
}
