package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
)

type IssueTokensService struct {
	jwtToken     jwt.JwtTokenGenerate
	cfg          config.Tokens
	refreshToken *RefreshToken
	accessToken  *AccessToken
}

func NewIssueTokensService(jwtToken jwt.JwtTokenGenerate, cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		refreshToken: NewRefreshToken(cfg.RefreshToken),
		accessToken:  NewAccessToken(cfg.AccessToken),
		jwtToken:     jwtToken,
		cfg:          cfg,
	}
}

func (t *IssueTokensService) GenerateAccessAndRefreshToken(data models.TokenPayload) (*models.AuthTokens, error) {
	refreshToken, err := t.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new refresh token error: %v", err)
	}

	accessToken, err := t.GenerateAccessToken(refreshToken, data)
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new access token error: %v", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *IssueTokensService) GenerateAccessToken(refreshToken string, data models.TokenPayload) (string, error) {
	accessToken, err := t.accessToken.generateNewAccessToken(refreshToken, data)
	if err != nil {
		return "", fmt.Errorf("generateTokensPair: generating new access token error: %v", err)
	}
	return accessToken, nil
}

func (t *IssueTokensService) GenerateRefreshToken() (string, error) {
	refreshToken, err := t.refreshToken.generateNewRefreshToken()
	if err != nil {
		return "", fmt.Errorf("generateTokensPair: generating new access token error: %v", err)
	}
	return refreshToken, nil
}

func (t *IssueTokensService) ParseRefreshToken(refreshToken string) error {
	return t.refreshToken.parseRefreshToken(refreshToken)
}

func (t *IssueTokensService) ParseAccessToken(tokens models.AuthTokens) (*models.AccessTokenInfo, error) {
	return t.accessToken.parseAccessToken(tokens)
}
