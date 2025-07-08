package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthService struct {
	db AuthStrore
}

func NewAuthService(db AuthStrore) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (us *AuthService) GetTokens(ctx context.Context, data models.DataFromRequestGetTokens) (*models.AuthTokens, error) {
	// проверить сущ-ет ли user в БД

	// Выпустить токен

	// Сохранить инфо о входе в БД

	createdGood, err := us.db.CreateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGoods - us.db.CreateGoods: %w", err)
	}

	return createdGood, nil
}

func (us *AuthService) UpdateGood(ctx context.Context, data models.DataFromRequestGoodsUpdate) (*models.Goods, error) {
	updGood, err := us.db.UpdateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - UpdateGood - us.db.UpdateGoo: %w", err)
	}

	return updGood, nil
}

func (us *AuthService) DeleteGood(ctx context.Context, data models.DataFromRequestGoodsDelete) (*models.Goods, error) {
	delData, err := us.db.DeleteGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - DeleteGood - us.db.DeleteGood: %w", err)
	}
	return delData, nil
}
