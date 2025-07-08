package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token"
	"golang.org/x/crypto/bcrypt"
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
	user, err := us.db.GetUser(ctx, data.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.GetUser error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("GetTokens - user info is nil")
	}

	// Выпустить токены
	tokens, err := us.tokensGenerate.GenerateTokensPair(data.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GenerateTokensPair error: %w", err)
	}

	// хешируем refresh токен
	hashRefreshToken, err := us.bcryptRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - bcryptRefreshToken error: %w", err)
	}

	// Сохранить инфо о входе в БД
	if err = us.db.SaveUserLogin(ctx, models.SaveLoginDataDbRequest{
		UserID:       data.UserID,
		RefreshToken: hashRefreshToken,
		UserAgent:    data.UserAgent,
		IpAddress:    data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return tokens, nil
}

func (us *AuthService) bcryptRefreshToken(refreshToken string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
