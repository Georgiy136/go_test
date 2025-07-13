package service

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/client"
	"github.com/Georgiy136/go_test/auth_service/helpers"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/app_errors"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/repo"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate"
	"github.com/Georgiy136/go_test/auth_service/internal/templates"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthService struct {
	notificationClient *client.NotificationClient
	issueTokensService *token_generate.IssueTokensService
	crypter            *crypter.Crypter
	db                 repo.AuthDBStore
}

func NewAuthService(
	issueTokensService *token_generate.IssueTokensService,
	crypter *crypter.Crypter,
	notificationClient *client.NotificationClient,
	db repo.AuthDBStore,
) *AuthService {
	return &AuthService{
		issueTokensService: issueTokensService,
		notificationClient: notificationClient,
		crypter:            crypter,
		db:                 db,
	}
}

func (us *AuthService) GetTokens(ctx context.Context, data models.DataFromRequestGetTokens) (*models.AuthTokens, error) {
	// генерируем id сессии
	sessionID := uuid.New().String()

	refreshToken, err := us.issueTokensService.RefreshToken.New()
	if err != nil {
		return nil, errors.Wrap(err, "RefreshToken.New error")
	}
	accessToken, err := us.issueTokensService.AccessToken.New(refreshToken, models.AccessTokenPayload{
		UserID:    data.UserID,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "AccessToken.New error")
	}

	refreshTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "a.crypter.Encrypt refreshToken error")
	}
	accessTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "a.crypter.Encrypt accessToken error")
	}

	if err = us.db.SaveUserSession(ctx, models.LoginInfo{
		UserID:    data.UserID,
		SessionID: sessionID,
		Token:     helpers.HashSha256(refreshTokenEncrypted),
		UserAgent: data.UserAgent,
		IpAddress: data.IpAddress,
	}); err != nil {
		return nil, errors.Wrap(err, "GetTokens - us.db.SaveUserLogin error")
	}

	return &models.AuthTokens{
		AccessToken:  accessTokenEncrypted,
		RefreshToken: refreshTokenEncrypted,
	}, nil
}

func (us *AuthService) UpdateTokens(ctx context.Context, data models.DataFromRequestUpdateTokens) (*models.AuthTokens, error) {
	refreshTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error")
	}
	accessTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt accessToken error")
	}

	if err = us.issueTokensService.RefreshToken.Parse(refreshTokenDecoded); err != nil {
		if !errors.Is(err, app_errors.TokenIsExpiredError) {
			return nil, errors.Wrap(err, "UpdateTokens - ParseRefreshToken error")
		}
	}

	accessTokenInfo, err := us.issueTokensService.AccessToken.Parse(models.AuthTokens{
		AccessToken:  accessTokenDecoded,
		RefreshToken: refreshTokenDecoded,
	})
	if err != nil {
		if !errors.Is(err, app_errors.TokenIsExpiredError) {
			return nil, errors.Wrap(err, "UpdateTokens - ParseAccessToken error")
		}
	}

	// ищем инфо о входе в БД по user_id и session_id
	loginInfo, err := us.db.GetUserSession(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			return nil, app_errors.SessionUserNotFoundError
		default:
			return nil, errors.Wrap(err, "UpdateTokens - us.db.GetUserSession error")
		}
	}

	// Сверяем совпадают ли refresh токен с хешированным в БД
	if !strings.EqualFold(helpers.HashSha256(data.RefreshToken), loginInfo.Token) {
		return nil, errors.Wrap(err, "UpdateTokens - RefreshToken does not match in db")
	}

	// Сверяем совпадают ли User-Agent
	if !strings.EqualFold(data.UserAgent, loginInfo.UserAgent) {
		go func() {
			if err = us.Logout(ctx, models.DataFromRequestLogout{AccessToken: data.AccessToken, RefreshToken: data.RefreshToken}); err != nil { // деавторизуем пользователя
				logrus.Errorf("UserAgent not match in db, failed to logout: %v", err)
			}
		}()

		return nil, app_errors.UserAgentNotMatchInDB
	}

	// Сверяем совпадает ли ip-адрес
	if !strings.EqualFold(data.IpAddress, loginInfo.IpAddress) {
		// уведомляем пользователя о входе с нового ip адреса
		go func() {
			if err = us.notificationClient.SendNewSignInNotification(accessTokenInfo.UserID, fmt.Sprintf(templates.NewSignInNotificationCommonMsg, data.IpAddress, data.UserAgent)); err != nil {
				logrus.Errorf("UpdateTokens - SendNewSignInNotification error: %v", err)
			}
		}()
	}

	// удаляем старую сессию в БД
	go func() {
		if err = us.db.DeleteUserSession(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID); err != nil {
			logrus.Errorf("UpdateTokens - DeleteUserSession error: %v", err)
		}
	}()

	// выпускаем новые токены
	return us.GetTokens(ctx, models.DataFromRequestGetTokens{
		UserID:    accessTokenInfo.UserID,
		UserAgent: data.UserAgent,
		IpAddress: data.IpAddress,
	})
}

func (us *AuthService) GetUser(ctx context.Context, data models.DataFromRequestGetUser) (*models.User, error) {
	refreshTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error")
	}
	accessTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt accessToken error")
	}

	if err = us.issueTokensService.RefreshToken.Parse(refreshTokenDecoded); err != nil {
		switch {
		case errors.Is(err, app_errors.TokenIsExpiredError):
			return nil, app_errors.TokenIsExpiredError
		default:
			return nil, errors.Wrap(err, "RefreshToken Parse error")
		}
	}

	accessTokenInfo, err := us.issueTokensService.AccessToken.Parse(models.AuthTokens{
		AccessToken:  accessTokenDecoded,
		RefreshToken: refreshTokenDecoded,
	})
	if err != nil {
		switch {
		case errors.Is(err, app_errors.TokenIsExpiredError):
			return nil, app_errors.TokenIsExpiredError
		default:
			return nil, errors.Wrap(err, "AccessToken Parse error")
		}
	}

	// проверяем инфо о входе в БД по user_id и session_id
	loginInfo, err := us.db.GetUserSession(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			return nil, app_errors.SessionUserNotFoundError
		default:
			return nil, errors.Wrap(err, "UpdateTokens - us.db.GetUserSession error")
		}
	}
	// Сверяем совпадают ли refresh токен с хешированным в БД
	if !strings.EqualFold(helpers.HashSha256(data.RefreshToken), loginInfo.Token) {
		return nil, errors.Wrap(err, "UpdateTokens - RefreshToken does not match in db")
	}

	return &models.User{UserID: accessTokenInfo.UserID}, nil
}

func (us *AuthService) Logout(ctx context.Context, data models.DataFromRequestLogout) error {
	refreshTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.RefreshToken)
	if err != nil {
		return errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error")
	}
	accessTokenDecoded, err := us.crypter.DecodeFromBase64AndDecrypt(data.AccessToken)
	if err != nil {
		return errors.Wrap(err, "UpdateTokens - DecodeFromBase64AndDecrypt accessToken error")
	}

	if err = us.issueTokensService.RefreshToken.Parse(refreshTokenDecoded); err != nil {
		if !errors.Is(err, app_errors.TokenIsExpiredError) {
			return errors.Wrap(err, "RefreshToken Parse error")
		}
	}
	accessTokenInfo, err := us.issueTokensService.AccessToken.Parse(models.AuthTokens{
		AccessToken:  accessTokenDecoded,
		RefreshToken: refreshTokenDecoded,
	})
	if err != nil {
		if !errors.Is(err, app_errors.TokenIsExpiredError) {
			return errors.Wrap(err, "AccessToken Parse error")
		}
	}

	// проверяем инфо о входе в БД по user_id и session_id
	loginInfo, err := us.db.GetUserSession(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			return app_errors.SessionUserNotFoundError
		default:
			return errors.Wrap(err, "Logout - us.db.GetUserSession error")
		}
	}
	// Сверяем совпадают ли refresh токен с хешированным в БД
	if !strings.EqualFold(helpers.HashSha256(data.RefreshToken), loginInfo.Token) {
		return errors.Wrap(err, "UpdateTokens - RefreshToken does not match in db")
	}

	// удаляем старую сессию в БД
	go func() {
		if err = us.db.DeleteUserSession(ctx, accessTokenInfo.UserID, accessTokenInfo.SessionID); err != nil {
			logrus.Errorf("UpdateTokens - DeleteUserSession error: %v", err)
		}
	}()

	return nil
}
