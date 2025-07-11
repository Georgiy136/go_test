package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/constant"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"time"
)

type IssueTokensService struct {
	tokenGenerator jwt.JwtTokenGenerate
	cfg            config.Tokens
	crypter        *crypter.Crypter
	refreshToken   *RefreshToken
	accessToken    *AccessToken
}

func NewIssueTokensService(jwtToken jwt.JwtTokenGenerate, crypter *crypter.Crypter, cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		refreshToken:   NewRefreshToken(jwtToken, cfg.RefreshToken),
		accessToken:    NewAccessToken(jwtToken, cfg.AccessToken),
		crypter:        crypter,
		tokenGenerator: jwtToken,
		cfg:            cfg,
	}
}

func (t *IssueTokensService) GenerateToken(tokenType constant.TokenType) (string, error) {
	switch tokenType {
	case constant.RefreshToken:
		return t.tokenGenerator.GenerateToken(t.cfg.RefreshToken.SignedKey, 1*time.Hour, "")
	case constant.AccessToken:
		return t.tokenGenerator.GenerateToken(t.cfg.AccessToken.SignedKey, 1*time.Hour, "")
	default:
		return "", fmt.Errorf("invalid token type")
	}

	// закодировать токены
	refreshTokenEncrypted, err := us.crypter.EncryptAndEncodeToBase64(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt refreshToken error: %w", err)
	}
}
