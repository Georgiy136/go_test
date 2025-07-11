package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/client"
	"github.com/Georgiy136/go_test/auth_service/helpers"
	"github.com/Georgiy136/go_test/auth_service/internal/common"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/app_errors"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate"
	"github.com/google/uuid"
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
	if _, err := us.getUserInfoClient.GetUserInfo(ctx, data.UserID); err != nil {
		return nil, fmt.Errorf("UpdateTokens - GetUserInfo error: %w", err)
	}

	// генерируем id сессии
	sessionID := uuid.New().String()

	refreshToken, err := us.issueTokensService.RefreshToken.New()
	if err != nil {
		return nil, fmt.Errorf("RefreshToken.New error: %w", err)
	}
	accessToken, err := us.issueTokensService.AccessToken.New(refreshToken, models.AccessTokenPayload{
		UserID:    data.UserID,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("AccessToken.New error: %w", err)
	}

	refreshTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt refreshToken error: %w", err)
	}
	accessTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(accessToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
	}

	logrus.Infof("helpers.HashSha256(refreshTokenEncrypted): %s", helpers.HashSha256(refreshTokenEncrypted))

	if err = us.db.SaveUserLogin(ctx, models.LoginInfo{
		UserID:       data.UserID,
		SessionID:    sessionID,
		RefreshToken: helpers.HashSha256(refreshTokenEncrypted),
		UserAgent:    data.UserAgent,
		IpAddress:    data.IpAddress,
	}); err != nil {
		return nil, fmt.Errorf("GetTokens - us.db.SaveUserLogin error: %w", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessTokenEncrypted,
		RefreshToken: refreshTokenEncrypted,
	}, nil
}

func (us *AuthService) UpdateTokens(ctx context.Context, data models.DataFromRequestUpdateTokens) (*models.AuthTokens, error) {
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

	if err = us.issueTokensService.RefreshToken.Parse(refreshTokenDecoded); err != nil {
		switch {
		case errors.Is(err, app_errors.TokenIsExpiredError):
			refreshTokenIsExpired = true
		default:
			return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
		}
	}

	accessTokenInfo, err := us.issueTokensService.AccessToken.Parse(models.AuthTokens{
		AccessToken:  accessTokenDecoded,
		RefreshToken: refreshTokenDecoded,
	})
	if err != nil {
		switch {
		case errors.Is(err, app_errors.TokenIsExpiredError):
			accessTokenIsExpired = true
		default:
			return nil, fmt.Errorf("UpdateTokens - ParseRefreshToken error: %w", err)
		}
	}

	// Проверяем сущ-ет ли пользователь
	if _, err = us.getUserInfoClient.GetUserInfo(ctx, accessTokenInfo.UserID); err != nil {
		return nil, fmt.Errorf("UpdateTokens - GetUserInfo error: %w", err)
	}

	// ищем инфо о входе в БД по user_id и session_id
	loginInfo, err := us.db.GetUserSignIn(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - us.db.GetUserSignIn error: %w", err)
	}

	// Сверяем совпадают ли refresh токен с захешированным в БД
	if !strings.EqualFold(helpers.HashSha256(data.RefreshToken), loginInfo.RefreshToken) {
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
		newAccessToken, err := us.issueTokensService.AccessToken.New(refreshTokenDecoded, models.AccessTokenPayload{
			UserID:    accessTokenInfo.UserID,
			SessionID: accessTokenInfo.SessionID,
		})
		if err != nil {
			return nil, fmt.Errorf("UpdateTokens - us.tokensGenerate.GenerateNewAccessToken error: %w", err)
		}
		newAccessTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(newAccessToken)
		if err != nil {
			return nil, fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
		}

		return &models.AuthTokens{
			AccessToken:  newAccessTokenEncrypted,
			RefreshToken: data.RefreshToken,
		}, nil
	}

	return &models.AuthTokens{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}, nil
}
