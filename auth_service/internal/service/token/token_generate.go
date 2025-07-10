package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type IssueTokensService struct {
	refreshToken *RefreshToken
	accessToken  *AccessToken
}

func NewIssueTokensService(cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		refreshToken: NewRefreshToken(cfg.RefreshToken),
		accessToken:  NewAccessToken(cfg.AccessToken),
	}
}

func (t *IssueTokensService) GenerateTokensPair(data models.TokenPayload) (*models.AuthTokens, error) {
	refreshToken, err := t.refreshToken.generateNewRefreshToken(data.RefreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new refresh token error: %v", err)
	}

	accessToken, err := t.accessToken.generateNewAccessToken(refreshToken, data)
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new access token error: %v", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *IssueTokensService) ParseRefreshToken(refreshToken string) (*models.RefreshTokenInfo, error) {
	return t.refreshToken.parseRefreshToken(refreshToken)
}

func (t *IssueTokensService) ParseAccessToken(tokens models.AuthTokens) (*models.AccessTokenInfo, error) {
	return t.accessToken.parseAccessToken(tokens.AccessToken, tokens.RefreshToken)
}
