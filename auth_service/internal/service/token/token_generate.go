package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"strconv"
)

type IssueTokensService struct {
	refreshToken *RefreshToken
	accessToken  *AccessToken
}

func NewIssueTokensService(cfg *config.Tokens) IssueTokensStore {
	return &IssueTokensService{
		refreshToken: NewRefreshToken(cfg.RefreshToken),
		accessToken:  NewAccessToken(cfg.AccessToken),
	}
}

func (t *IssueTokensService) GenerateTokensPair(userID int) (*models.AuthTokens, error) {
	// сгенерить refresh токен
	refreshToken, err := t.refreshToken.generateNewRefreshToken(strconv.Itoa(userID))
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new refresh token error: %v", err)
	}

	// сгенерить access токен
	accessToken, err := t.accessToken.generateNewAccessToken(refreshToken, strconv.Itoa(userID))
	if err != nil {
		return nil, fmt.Errorf("generateTokensPair: generating new refresh token error: %v", err)
	}

	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
