package token_generate

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/tokens"
)

type IssueTokensService struct {
	crypter      *crypter.Crypter
	refreshToken tokens.RefreshTokenStore
	accessToken  tokens.AccessTokenStore
}

func NewIssueTokensService(
	refreshToken tokens.RefreshTokenStore,
	accessToken tokens.AccessTokenStore,
	crypter *crypter.Crypter) *IssueTokensService {
	return &IssueTokensService{
		refreshToken: refreshToken,
		accessToken:  accessToken,
		crypter:      crypter,
	}
}

func (t *IssueTokensService) GenerateRefreshAndAccessTokens(payload models.AccessTokenPayload) (*models.AuthTokens, error) {
	refreshToken, err := t.refreshToken.New()
	if err != nil {
		return nil, fmt.Errorf("GeneratePairsRefreshAndAccessTokens - NewRefreshToken error: %w", err)
	}

	accessToken, err := t.accessToken.New(refreshToken, payload)
	if err != nil {
		return nil, fmt.Errorf("GeneratePairsRefreshAndAccessTokens - NewAccessToken error: %w", err)
	}
	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *IssueTokensService) ParseRefreshToken(refreshToken string) error {
	return t.refreshToken.Parse(refreshToken)
}

func (t *IssueTokensService) ParseAccessToken(tokensPair models.AuthTokens) (*models.AccessTokenPayload, error) {
	return t.accessToken.Parse(tokensPair)
}

func (t *IssueTokensService) NewAccessToken(refreshToken string, payload models.AccessTokenPayload) (string, error) {
	return t.accessToken.New(refreshToken, payload)
}
