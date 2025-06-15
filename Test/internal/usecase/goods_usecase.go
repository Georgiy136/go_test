package usecase

import (
	"context"
	"fmt"
	"myapp/internal/models"

	"github.com/google/uuid"
)

type GoodsUseCases struct {
	store GoodsStrore
}

func NewGoodsUsecases(st GoodsStrore) *GoodsUseCases {
	return &GoodsUseCases{
		store: st,
	}
}

func (us *GoodsUseCases) AddGood(ctx context.Context, p models.Goods) (*models.Goods, error) {
	err := us.store.CreateGoods(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGood - us.store.CreateGood: %w", err)
	}
	return &p, nil
}

func (us *GoodsUseCases) GetAllGoods(ctx context.Context) ([]models.Goods, error) {
	p, err := us.store.GetAllGoods(ctx)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetAllGoods - us.store.GetAllGoods: %w", err)
	}
	return p, nil
}

func (us *GoodsUseCases) DeleteGood(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("GoodUseCases - DeleteGood - uuid.Parse: %w", err)
	}
	err = us.store.DeleteGoods(ctx, uid)
	if err != nil {
		return fmt.Errorf("GoodUseCases - DeleteGood - us.store.DeleteGood: %w", err)
	}
	return nil
}

func (us *GoodsUseCases) UpdateGood(ctx context.Context, id string, p models.Goods) (*models.Goods, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - uuid.Parse: %w", err)
	}
	Good, err := us.store.UpdateGoods(ctx, uid, p)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - us.store.UpdateGood: %w", err)
	}
	return Good, nil
}

func (us *GoodsUseCases) GetOneGood(ctx context.Context, id string) (*models.Goods, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetOneGood - uuid.Parse: %w", err)
	}
	p, err := us.store.GetOneGoods(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - GetOneGood - us.store.GetOneGood: %w", err)
	}
	return p, nil
}
