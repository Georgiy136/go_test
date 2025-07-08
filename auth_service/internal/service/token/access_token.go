package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
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

	tokenString, err := a.crypter.EncryptAndEncodeToBase64(jwtToken)
	if err != nil {
		return "", fmt.Errorf("a.crypter.Encrypt error: %w", err)
	}

	return tokenString, nil
}

func (a *AccessToken) decodeAccessToken(accessToken, refreshToken string) (*models.TokenInfo, error) {
	accessTokenDecode, err := a.crypter.DecodeFromBase64AndDecrypt(accessToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt error: %w", err)
	}

	refreshTokenDecode, err := a.crypter.DecodeFromBase64AndDecrypt(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("a.crypter.Encrypt error: %w", err)
	}

	signedString := a.getSignedString(refreshTokenDecode)

	token, err := jwt.ParseWithClaims(accessTokenDecode, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signedString), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("decodeAccessToken jwt.Parse error: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return &models.TokenInfo{
			Issuer:    claims.Issuer,
			Payload:   claims.Subject,
			ExpiredAt: claims.ExpiresAt.Time,
			IssuedAt:  claims.ExpiresAt.Time,
		}, nil
	}

	return nil, fmt.Errorf("decodeAccessToken error: %w", err)
}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}
