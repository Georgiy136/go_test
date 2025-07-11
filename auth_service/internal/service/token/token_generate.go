package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/constant"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"time"
)

type IssueTokensService struct {
	jwtToken     jwt.JwtTokenGenerate
	cfg          config.Tokens
	refreshToken *RefreshToken
	accessToken  *AccessToken
}

func NewIssueTokensService(jwtToken jwt.JwtTokenGenerate, cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		refreshToken: NewRefreshToken(jwtToken, cfg.RefreshToken),
		accessToken:  NewAccessToken(jwtToken, cfg.AccessToken),
		jwtToken:     jwtToken,
		cfg:          cfg,
	}
}

func (t *IssueTokensService) GenerateToken(tokenType constant.TokenType) (string, error) {
	switch tokenType {
	case constant.RefreshToken:
		return t.jwtToken.GenerateToken(t.cfg.RefreshToken.SignedKey, 1*time.Hour, "")
	case constant.AccessToken:
		return t.jwtToken.GenerateToken(t.cfg.AccessToken.SignedKey, 1*time.Hour, "")
	default:
		return "", fmt.Errorf("invalid token type")
	}
}
