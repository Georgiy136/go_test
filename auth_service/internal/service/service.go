package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/helpers"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token"
	"strings"
)

type AuthService struct {
	db             AuthStore
	tokensGenerate token.IssueTokensService
}

func NewAuthService(tokensGenerate token.IssueTokensService, db AuthStore) *AuthService {
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
	if err = us.db.SaveUserLogin(ctx, models.LoginInfo{
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
	refreshToken, err := us.tokensGenerate.DecodeFromBase64AndDecrypt(data.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error: %w", err)
	}
	accessToken, err := us.tokensGenerate.DecodeFromBase64AndDecrypt(data.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - DecodeFromBase64AndDecrypt accessToken error: %w", err)
	}

	refreshTokenInfo, err := us.tokensGenerate.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
	}
	accessTokenInfo, err := us.tokensGenerate.ParseAccessToken(accessToken, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
	}

	// сверяем RefreshTokenID
	if refreshTokenInfo.RefreshTokenID != accessTokenInfo.RefreshTokenID {
		return nil, fmt.Errorf("UpdateTokens - RefreshTokenID does not match")
	}

	// Проверяем есть ли пользователь в БД
	user, err := us.db.GetUser(ctx, accessTokenInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - us.db.GetUser error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("UpdateTokens - GetUser - user info is nil")
	}

	// ищем инфо о входе по refresh токену в БД
	loginInfo, err := us.db.GetSignInByRefreshTokenID(ctx, accessTokenInfo.RefreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - us.db.GetSignInByRefreshTokenID error: %w", err)
	}

	// Сверяем совпадают ли refresh токен с захешированным в БД
	if !strings.EqualFold(helpers.HashSha512(data.RefreshToken), loginInfo.RefreshToken) {
		return nil, fmt.Errorf("UpdateTokens - RefreshToken does not match in db")
	}

	// Сверяем совпадают ли User-Agent
	if !strings.EqualFold(data.UserAgent, loginInfo.UserAgent) {
		return nil, fmt.Errorf("UpdateTokens - User-Agent does not match in db")
	}

	// Сверяем совпадает ли ip-адрес
	if !strings.EqualFold(data.IpAddress, loginInfo.IpAddress) {
		// уведомляем пользователя о входе с нового ip адреса
		// ...
	}

	// Сохранить инфо о входе в БД
	if err = us.db.SaveUserLogin(ctx, models.SaveLoginInfoDbRequest{
		UserID:       data.UserID,
		RefreshToken: helpers.HashSha512(tokens.RefreshToken),
		UserAgent:    data.UserAgent,
		IpAddress:    data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return tokens, nil
}
