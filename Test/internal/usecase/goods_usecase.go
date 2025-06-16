package usecase

import (
	"context"
	"fmt"
	"myapp/internal/models"

	"github.com/google/uuid"
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

func (us *GoodsUseCases) AddGoods(ctx context.Context, p models.Goods) (*models.Goods, error) {
	err := us.db.CreateGoods(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGood - us.db.CreateGood: %w", err)
	}
	return &p, nil
}

func (us *GoodsUseCases) GetAllGoods(ctx context.Context) ([]models.Goods, error) {
	p, err := us.db.GetAllGoods(ctx)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetAllGoods - us.db.GetAllGoods: %w", err)
	}
	return p, nil
}

func (us *GoodsUseCases) DeleteGood(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("GoodUseCases - DeleteGood - uuid.Parse: %w", err)
	}
	err = us.db.DeleteGoods(ctx, uid)
	if err != nil {
		return fmt.Errorf("GoodUseCases - DeleteGood - us.db.DeleteGood: %w", err)
	}
	return nil
}

func (us *GoodsUseCases) UpdateGood(ctx context.Context, id string, p models.Goods) (*models.Goods, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - uuid.Parse: %w", err)
	}
	Good, err := us.db.UpdateGoods(ctx, uid, p)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - us.db.UpdateGood: %w", err)
	}
	return Good, nil
}

func (us *GoodsUseCases) GetOneGood(ctx context.Context, id string) (*models.Goods, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetOneGood - uuid.Parse: %w", err)
	}
	p, err := us.db.GetOneGoods(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetOneGood - us.db.GetOneGood: %w", err)
	}
	return p, nil
}
