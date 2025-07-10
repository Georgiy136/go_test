package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"time"
)

type AccessToken struct {
	cfg           config.AccessToken
	tokenLifetime time.Duration
}

func NewAccessToken(cfg config.AccessToken) *AccessToken {
	tokenLifetime, err := time.ParseDuration(cfg.TokenLifetime)
	if err != nil {
		logrus.Fatalf("NewAccessToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &AccessToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
	}
}

func (a *AccessToken) generateNewAccessToken(refreshToken string, accessTokenPayload models.TokenPayload) (string, error) {
	payloadBytes, err := jsoniter.MarshalToString(accessTokenPayload)
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken: json marshal payload err: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenLifetime)),
		Subject:   payloadBytes,
	})

	jwtToken, err := token.SignedString([]byte(a.getSignedString(refreshToken)))
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken token.SignedString error: %w", err)
	}

	return jwtToken, nil
}

func (a *AccessToken) parseAccessToken(tokens models.AuthTokens) (*models.AccessTokenInfo, error) {
	token, err := jwt.ParseWithClaims(tokens.AccessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.getSignedString(tokens.RefreshToken)), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("decodeAccessToken jwt.Parse error: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		payload, err := a.getAccessTokenPayload(claims.Subject)
		if err != nil {
			return nil, fmt.Errorf("decodeAccessTokenv getAccessTokenPayload error: %w", err)
		}

		return &models.AccessTokenInfo{
			Issuer:         claims.Issuer,
			UserID:         payload.UserID,
			RefreshTokenID: payload.RefreshTokenID,
			ExpiredAt:      claims.ExpiresAt.Time,
			IssuedAt:       claims.IssuedAt.Time,
		}, nil
	}

	return nil, fmt.Errorf("decodeAccessToken error: %w", err)
}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}

func (a *AccessToken) getAccessTokenPayload(payload string) (*models.TokenPayload, error) {
	var payloadData models.TokenPayload
	if err := jsoniter.UnmarshalFromString(payload, &payloadData); err != nil {
		return nil, err
	}

	return &payloadData, nil
}
