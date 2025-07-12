package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
)

type IssueTokensService struct {
	tokenGenerator jwt.JwtTokenGenerate
	crypter        *crypter.Crypter
	refreshToken   *refreshToken
	accessToken    *accessToken
}

func NewIssueTokensService(jwtTokenGenerator jwt.JwtTokenGenerate, crypter *crypter.Crypter, cfg config.Tokens) *IssueTokensService {
	return &IssueTokensService{
		tokenGenerator: jwtTokenGenerator,
		refreshToken:   NewRefreshToken(jwtTokenGenerator, cfg.RefreshToken),
		accessToken:    NewAccessToken(jwtTokenGenerator, cfg.AccessToken),
		crypter:        crypter,
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
