package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type IssueTokensService struct {
	refreshToken *RefreshToken
	accessToken  *AccessToken
	crypter      *crypter
}

func NewIssueTokensService(cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		refreshToken: NewRefreshToken(cfg.RefreshToken),
		accessToken:  NewAccessToken(cfg.AccessToken),
		crypter:      NewCrypter(cfg.Crypter.SignedKey),
	}
}

func (t *IssueTokensService) GenerateTokensPair(data models.TokenPayload) (*models.AuthTokens, error) {
	// генерим refresh токен
	refreshToken, err := t.refreshToken.generateNewRefreshToken(data.RefreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new refresh token error: %v", err)
	}

	// генерим access токен
	accessToken, err := t.accessToken.generateNewAccessToken(refreshToken, data)
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new access token error: %v", err)
	}

	// доп-но зашифровываем
	refreshTokenEncrypted, err := t.crypter.EncryptAndEncodeToBase64(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt refreshToken error: %w", err)
	}
	accessTokenEncrypted, err := t.crypter.EncryptAndEncodeToBase64(accessToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessTokenEncrypted,
		RefreshToken: refreshTokenEncrypted,
	}, nil
}

func (t *IssueTokensService) DecodeFromBase64AndDecrypt(data string) (string, error) {
	return t.crypter.DecodeFromBase64AndDecrypt(data)
}

func (t *IssueTokensService) ParseRefreshToken(refreshToken string) (*models.RefreshTokenInfo, error) {
	return t.refreshToken.parseRefreshToken(refreshToken)
}

func (t *IssueTokensService) ParseAccessToken(accessToken, refreshToken string) (*models.AccessTokenInfo, error) {
	return t.accessToken.parseAccessToken(accessToken, refreshToken)
}
