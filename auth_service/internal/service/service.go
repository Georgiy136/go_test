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
	if err := us.GetUserDBReq(ctx, data.UserID); err != nil {
		return nil, fmt.Errorf("GetTokens - GetUserDBReq error: %w", err)
	}

	// Получаем уник-ый refresh_token_id
	refreshTokenID, err := us.db.GetRefreshTokenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GetRefreshTokenID error: %w", err)
	}

	// Выпустить токены
	tokens, err := us.tokensGenerate.GenerateTokensPair(refreshTokenID, data.UserID)
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

func (us *AuthService) GetUserDBReq(ctx context.Context, userID int) error {
	user, err := us.db.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("GetUserDBReq - us.db.GetUser error: %w", err)
	}
	if user == nil {
		return fmt.Errorf("GetUserDBReq - user info is nil")
	}
	return nil
}
