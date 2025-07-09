package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/helpers"
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
	// Проверяем есть ли пользователь в БД
	user, err := us.db.GetUser(ctx, data.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.GetUser error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("GetTokens - GetUser - user info is nil")
	}

	// Получаем уникальный refresh_token_id из БД (сдвигаем сиквенс)
	refreshTokenID, err := us.db.GetRefreshTokenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GetRefreshTokenID error: %w", err)
	}

	// Выпустить токены
	tokens, err := us.tokensGenerate.GenerateTokensPair(models.TokenPayload{
		UserID:         data.UserID,
		RefreshTokenID: refreshTokenID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GenerateTokensPair error: %w", err)
	}

	// Сохранить инфо о входе в БД
	if err = us.db.SaveUserLogin(ctx, models.SaveLoginDataDbRequest{
		UserID:         data.UserID,
		RefreshTokenID: refreshTokenID,
		RefreshToken:   helpers.HashSha512(tokens.RefreshToken),
		UserAgent:      data.UserAgent,
		IpAddress:      data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return tokens, nil
}

func (us *AuthService) UpdateTokens(ctx context.Context, data models.DataFromRequestUpdateTokens) (*models.AuthTokens, error) {

	// Выпустить токены
	tokens, err := us.tokensGenerate.GenerateTokensPair(data.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GenerateTokensPair error: %w", err)
	}

	// Сохранить инфо о входе в БД
	if err = us.db.SaveUserLogin(ctx, models.SaveLoginDataDbRequest{
		UserID:       data.UserID,
		RefreshToken: helpers.HashSha512(tokens.RefreshToken),
		UserAgent:    data.UserAgent,
		IpAddress:    data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return tokens, nil
}
