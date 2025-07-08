package service

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type AccessToken struct {
	cfg           config.AccessToken
	tokenLifetime time.Duration
	crypter       *crypter
}

func NewAccessToken(cfg config.AccessToken) *AccessToken {
	tokenLifetime, err := time.ParseDuration(cfg.SignedKey)
	if err != nil {
		logrus.Fatalf("NewAccessToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &AccessToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
		crypter:       NewCrypter(cfg.SignedKey),
	}
}

func (a *AccessToken) generateNewAccessToken(refreshToken string, payload string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenLifetime)),
		Subject:   payload,
	})

	signedString := a.getSignedString(refreshToken)

	jwtToken, err := token.SignedString([]byte(signedString))
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken token.SignedString error: %w", err)
	}

	tokenString, err := a.crypter.Encrypt(jwtToken)
	if err != nil {
		return "", fmt.Errorf("a.crypter.Encrypt error: %w", err)
	}

	return string(tokenString), nil
}

func (a *AccessToken) decodeAccessToken(accessToken, refreshToken string) (string, error) {
	signedString := a.getSignedString(refreshToken)

	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signedString), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", fmt.Errorf("decodeAccessToken jwt.Parse error: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("decodeAccessToken error: %w", err)
}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}
