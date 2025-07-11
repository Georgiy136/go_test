package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"github.com/go-faster/errors"
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

func (a *AccessToken) generateNewAccessToken(refreshToken string, accessTokenPayload models.AccessTokenPayload) (string, error) {
	return a.jwtToken.GenerateToken(a.getSignedString(refreshToken), a.tokenLifetime, accessTokenPayload)
}

func (a *AccessToken) parseAccessToken(tokens models.AuthTokens) (*models.AccessTokenPayload, error) {
	sub, err := a.jwtToken.ParseToken(tokens.AccessToken, a.getSignedString(tokens.RefreshToken))
	if err != nil {
		return nil, errors.Wrap(err, "parseAccessToken error")
	}

	payload, err := a.getAccessTokenPayload(sub)
	if err != nil {
		return nil, fmt.Errorf("decodeAccessToken a.getAccessTokenPayload error: %w", err)
	}

	return &models.AccessTokenPayload{
		UserID:         payload.UserID,
		RefreshTokenID: payload.RefreshTokenID,
	}, nil
}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}

func (a *AccessToken) getAccessTokenPayload(payload string) (*models.AccessTokenPayload, error) {
	var payloadData models.AccessTokenPayload
	if err := jsoniter.UnmarshalFromString(payload, &payloadData); err != nil {
		return nil, err
	}
	return &payloadData, nil
}
