package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token"
)

type AuthService struct {
	db             AuthStrore
	tokensGenerate token.IssueTokensStore
}

func NewAuthService(tokensGenerate token.IssueTokensStore, db AuthStrore) *AuthService {
	return &AuthService{
		db:             db,
		tokensGenerate: tokensGenerate,
	}
}

func (us *AuthService) GetTokens(ctx context.Context, data models.DataFromRequestGetTokens) (*models.AuthTokens, error) {
	// проверить сущ-ет ли user в БД
	// ...

	// Выпустить токены
	tokens, err := us.tokensGenerate.GenerateTokensPair(data.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GenerateTokensPair error: %w", err)
	}

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
