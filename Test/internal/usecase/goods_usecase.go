package usecase

import (
	"context"
	"fmt"
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
	goods, err := us.db.CreateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGood - us.db.CreateGoods: %w", err)
	}
	return goods, nil
}

func (us *GoodsUseCases) UpdateGood(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	Good, err := us.db.UpdateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - us.db.UpdateGood: %w", err)
	}
	return Good, nil
}

func (us *GoodsUseCases) DeleteGood(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	err := us.db.DeleteGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - DeleteGood - us.db.DeleteGood: %w", err)
	}
	return nil, nil
}

func (us *GoodsUseCases) ListGoods(ctx context.Context, data models.DataFromRequestGoodsList) (*models.GoodsListDBResponse, error) {
	p, err := us.db.ListGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetGoods - us.db.GetGoods: %w", err)
	}
	return p, nil
}

func (us *GoodsUseCases) ReprioritizeGood(ctx context.Context, data models.DataFromRequestReprioritizeGood) ([]models.Goods, error) {
	/*p, err := us.db.ListGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetGoods - us.db.GetGoods: %w", err)
	}*/
	return nil, nil
}
