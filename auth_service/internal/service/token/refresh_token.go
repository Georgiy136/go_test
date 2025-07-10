package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type RefreshToken struct {
	cfg           config.RefreshToken
	tokenLifetime time.Duration
}

func NewRefreshToken(cfg config.RefreshToken) *RefreshToken {
	tokenLifetime, err := time.ParseDuration(cfg.TokenLifetime)
	if err != nil {
		logrus.Fatalf("NewAccessToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &RefreshToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
	}
}

func (a *RefreshToken) generateNewRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenLifetime)),
	})

	jwtToken, err := token.SignedString([]byte(a.cfg.SignedKey))
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken token.SignedString error: %w", err)
	}

	return jwtToken, nil
}

func (a *RefreshToken) parseRefreshToken(refreshToken string) error {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.SignedKey), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return fmt.Errorf("decodeRefreshToken jwt.Parse error: %w", err)
	}

	if _, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return nil
	}

	return fmt.Errorf("decodeRefreshToken error: %w", err)
}
