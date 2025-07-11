package token

import (
	"errors"
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
	parsedAccessToken, err := jwt.Parse(tokens.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.getSignedString(tokens.RefreshToken)), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			err = jwt.ErrTokenExpired
		default:
			return nil, fmt.Errorf("decodeAccessToken jwt.Parse error: %w", err)
		}
	}

	payloadString, err := parsedAccessToken.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("decodeAccessToken claims.GetSubject error: %w", err)
	}

	payload, err := a.getAccessTokenPayload(payloadString)
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

func (a *AccessToken) getAccessTokenPayload(payload string) (*models.TokenPayload, error) {
	var payloadData models.TokenPayload
	if err := jsoniter.UnmarshalFromString(payload, &payloadData); err != nil {
		return nil, err
	}

	return &payloadData, nil
}
