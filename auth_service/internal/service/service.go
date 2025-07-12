package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/client"
	"github.com/Georgiy136/go_test/auth_service/helpers"
	"github.com/Georgiy136/go_test/auth_service/internal/common"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/jwt"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthService struct {
	notificationClient *client.NotificationClient
	getUserInfoClient  *client.UserInfoClient
	issueTokensService *token_generate.IssueTokensService
	crypter            *crypter.Crypter
	db                 AuthDBStore
}

func NewAuthService(
	issueTokensService *token_generate.IssueTokensService,
	crypter *crypter.Crypter,
	getUserInfoClient *client.UserInfoClient,
	notificationClient *client.NotificationClient,
	db AuthDBStore,
) *AuthService {
	return &AuthService{
		issueTokensService: issueTokensService,
		notificationClient: notificationClient,
		getUserInfoClient:  getUserInfoClient,
		crypter:            crypter,
		db:                 db,
	}
}

func (us *AuthService) GetTokens(ctx context.Context, data models.DataFromRequestGetTokens) (*models.AuthTokens, error) {
	// Проверяем сущ-ет ли пользователь
	//if _, err = us.getUserInfoClient.GetUserInfo(ctx, accessTokenInfo.UserID); err != nil {
	//	return nil, fmt.Errorf("UpdateTokens - GetUserInfo error: %w", err)
	//}

	// Получаем уникальный refresh_token_id из БД (сдвигаем сиквенс)
	refreshTokenID, err := us.db.GetRefreshTokenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GetRefreshTokenID error: %w", err)
	}

	// Выпустить токены
	tokens, err := us.issueTokensService.GenerateRefreshAndAccessTokens(models.AccessTokenPayload{
		UserID:         data.UserID,
		RefreshTokenID: refreshTokenID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetTokens - GenerateTokensPair error: %w", err)
	}

	// кодируем токены перед выпуском
	accessTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
	}
	refreshTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
	}

	//шифруем токен перед сохранением в БД
	hashRefreshToken := helpers.HashSha512(refreshTokenEncrypted)

	// Сохранить инфо о входе в БД
	if err = us.db.SaveUserLogin(ctx, models.LoginInfo{
		UserID:         data.UserID,
		RefreshTokenID: refreshTokenID,
		RefreshToken:   hashRefreshToken,
		UserAgent:      data.UserAgent,
		IpAddress:      data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessTokenEncrypted,
		RefreshToken: refreshTokenEncrypted,
	}, nil
}

func (us *AuthService) UpdateTokens(ctx context.Context, data models.DataFromRequestUpdateTokens) (*models.AuthTokens, error) {
	// декодируем токены
	refreshTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error: %w", err)
	}
	accessTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - DecodeFromBase64AndDecrypt accessToken error: %w", err)
	}

	var (
		refreshTokenIsExpired bool
		accessTokenIsExpired  bool
	)

	// парсим refresh токен
	if err = us.issueTokensService.ParseRefreshToken(refreshTokenDecoded); err != nil {
		switch {
		case errors.Is(err, jwt.TokenIsExpiredError):
			refreshTokenIsExpired = true
		default:
			return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
		}
	}

	// парсим access токен
	accessTokenInfo, err := us.issueTokensService.ParseAccessToken(models.AuthTokens{
		AccessToken:  accessTokenDecoded,
		RefreshToken: refreshTokenDecoded,
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.TokenIsExpiredError):
			accessTokenIsExpired = true
		default:
			return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
		}
	}

	// Проверяем сущ-ет ли пользователь
	//if _, err = us.getUserInfoClient.GetUserInfo(ctx, accessTokenInfo.UserID); err != nil {
	//	return nil, fmt.Errorf("UpdateTokens - GetUserInfo error: %w", err)
	//}

	// ищем инфо о входе в БД
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
		go func() {
			if err = us.notificationClient.SendNewSignInNotification(accessTokenInfo.UserID, fmt.Sprintf(common.NewSignInNotificationCommonMsg, data.IpAddress, data.UserAgent)); err != nil {
				logrus.Errorf("UpdateTokens - SendNewSignInNotification error: %v", err)
			}
		}()
	}

	if refreshTokenIsExpired {
		// удаляем старую сессию в БД
		// ...

		// выпускаем новые токены
		return us.GetTokens(ctx, models.DataFromRequestGetTokens{
			UserID:    accessTokenInfo.UserID,
			UserAgent: data.UserAgent,
			IpAddress: data.IpAddress,
		})
	}

	if accessTokenIsExpired {
		accessTokenDecoded, err = us.issueTokensService.NewAccessToken(refreshTokenDecoded, models.AccessTokenPayload{
			UserID:         accessTokenInfo.UserID,
			RefreshTokenID: accessTokenInfo.RefreshTokenID,
		})
		if err != nil {
			return nil, fmt.Errorf("UpdateTokens - us.tokensGenerate.GenerateNewAccessToken error: %w", err)
		}
	}

	return nil, fmt.Errorf("UpdateTokens - uknown err")
}
