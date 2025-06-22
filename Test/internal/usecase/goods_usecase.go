package usecase

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"myapp/internal/models"
)

type GoodsUseCases struct {
	db    GoodsStrore
	cache GoodsCache
}

func NewGoodsUsecases(db GoodsStrore, cache GoodsCache) *GoodsUseCases {
	return &GoodsUseCases{
		db:    db,
		cache: cache,
	}
}

func (us *GoodsUseCases) AddGoods(ctx context.Context, data models.DataFromRequestGoodsAdd) (*models.Goods, error) {
	createdGood, err := us.db.CreateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGoods - us.db.CreateGoods: %w", err)
	}

	return createdGood, nil
}

func (us *GoodsUseCases) UpdateGood(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	updGood, err := us.db.UpdateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - us.db.UpdateGoo: %w", err)
	}
	// очищаем из redis
	go func() {
		if err = us.cache.ClearGoods(ctx, data.GoodID, data.ProjectID); err != nil {
			logrus.Errorf("GoodUseCases - cache.ClearGoods: %v", err)
		}
	}()
	return updGood, nil
}

func (us *GoodsUseCases) DeleteGood(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	delData, err := us.db.DeleteGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - DeleteGood - us.db.DeleteGood: %w", err)
	}
	// очищаем из redis
	go func() {
		if err = us.cache.ClearGoods(ctx, data.GoodID, data.ProjectID); err != nil {
			logrus.Errorf("GoodUseCases - cache.ClearGoods: %v", err)
		}
	}()
	return delData, nil
}

func (us *GoodsUseCases) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsList, error) {
	// ищем запись в redis
	if data.GoodsID != nil && data.ProjectID != nil {
		goods, err := us.cache.GetGoods(ctx, *data.GoodsID, *data.ProjectID)
		if err != nil {
			logrus.Errorf("ListGoods - us.cache.GetGoods error: %v", err)
		}
		if goods != nil {
			return &models.GoodsList{Goods: []models.Goods{*goods}}, nil
		}
	}

	goodsList, err := us.db.ListGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - ListGoods - us.db.ListGoods: %w", err)
	}

	// сохраняем в redis запись
	if goodsList.Meta.Total == 1 {
		go func() {
			if err = us.cache.SaveGoods(ctx, *data.GoodsID, *data.ProjectID, goodsList.Goods[0]); err != nil {
				logrus.Errorf("ListGoods - us.cache.SaveGoods error: %v", err)
			}
		}()
	}

	return goodsList, nil
}

func (us *GoodsUseCases) ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) (*models.Goods, error) {
	goods, err := us.db.ReprioritizeGood(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetGoods - us.db.GetGoods: %w", err)
	}
	// очищаем из redis
	go func() {
		if err = us.cache.ClearGoods(ctx, data.GoodsID, data.ProjectID); err != nil {
			logrus.Errorf("GoodUseCases - cache.ClearGoods: %v", err)
		}
	}()
	return goods, nil
}
