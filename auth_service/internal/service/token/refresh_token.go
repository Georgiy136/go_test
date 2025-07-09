package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type RefreshToken struct {
	cfg           config.RefreshToken
	tokenLifetime time.Duration
	crypter       *crypter
}

func NewRefreshToken(cfg config.RefreshToken) *RefreshToken {
	tokenLifetime, err := time.ParseDuration(cfg.SignedKey)
	if err != nil {
		logrus.Fatalf("NewAccessToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &RefreshToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
		crypter:       NewCrypter(cfg.SignedKey),
	}
}

func (a *RefreshToken) generateNewRefreshToken(refreshTokenID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenLifetime)),
		ID:        fmt.Sprintf("%d", refreshTokenID),
	})

	jwtToken, err := token.SignedString([]byte(a.cfg.SignedKey))
	if err != nil {
		return "", fmt.Errorf("generateNewAccessToken token.SignedString error: %w", err)
	}

	tokenString, err := a.crypter.EncryptAndEncodeToBase64(jwtToken)
	if err != nil {
		return "", fmt.Errorf("a.crypter.Encrypt error: %w", err)
	}

	return tokenString, nil
}

func (a *RefreshToken) decodeRefreshToken(refreshToken string) (*models.RefreshTokenInfo, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.SignedKey), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, fmt.Errorf("decodeRefreshToken jwt.Parse error: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		tokenID, err := a.getAccessTokenID(claims.ID)
		if err != nil {
			return nil, fmt.Errorf("decodeRefreshToken jwt.getAccessTokenID error: %w", err)
		}

		return &models.RefreshTokenInfo{
			Issuer:         claims.Issuer,
			RefreshTokenID: tokenID,
			ExpiredAt:      claims.ExpiresAt.Time,
			IssuedAt:       claims.IssuedAt.Time,
		}, nil
	}

	return nil, fmt.Errorf("decodeRefreshToken error: %w", err)
}

func (a *RefreshToken) getAccessTokenID(id string) (int, error) {
	return strconv.Atoi(id)
}
