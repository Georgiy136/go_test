package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"time"
)

type AccessToken struct {
	jwtToken      jwt.JwtTokenGenerate
	cfg           config.AccessToken
	tokenLifetime time.Duration
}

func NewAccessToken(jwtToken jwt.JwtTokenGenerate, cfg config.AccessToken) *AccessToken {
	tokenLifetime, err := time.ParseDuration(cfg.TokenLifetime)
	if err != nil {
		logrus.Fatalf("NewAccessToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &AccessToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
		jwtToken:      jwtToken,
	}
}

func (a *AccessToken) generateNewAccessToken(refreshToken string, accessTokenPayload models.TokenPayload) (string, error) {
	payloadString, err := a.generateAccessTokenPayload(accessTokenPayload)
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken: json marshal payload err: %v", err)
	}

	return a.jwtToken.GenerateToken(a.getSignedString(refreshToken), a.tokenLifetime, payloadString)
}

func (a *AccessToken) parseAccessToken(tokens models.AuthTokens) (*models.AccessTokenInfo, error) {
	sub, err := a.jwtToken.ParseToken(tokens.AccessToken, a.getSignedString(tokens.RefreshToken))
	if err != nil {
		return nil, fmt.Errorf("generateNewAccessToken: json marshal payload err: %v", err)
	}

	payload, err := a.getAccessTokenPayload(sub)
	if err != nil {
		return nil, fmt.Errorf("decodeAccessToken a.getAccessTokenPayload error: %w", err)
	}

	return &models.AccessTokenInfo{
		UserID:         payload.UserID,
		RefreshTokenID: payload.RefreshTokenID,
	}, err

}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}

func (a *AccessToken) generateAccessTokenPayload(accessTokenPayload models.TokenPayload) (string, error) {
	payloadString, err := jsoniter.MarshalToString(accessTokenPayload)
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken: json marshal payload err: %v", err)
	}
	return payloadString, nil
}

func (a *AccessToken) getAccessTokenPayload(payload string) (*models.TokenPayload, error) {
	var payloadData models.TokenPayload
	if err := jsoniter.UnmarshalFromString(payload, &payloadData); err != nil {
		return nil, err
	}
	return &payloadData, nil
}
