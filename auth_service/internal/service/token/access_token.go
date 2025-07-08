package service

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type AccessToken struct {
	cfg     config.AccessToken
	crypter *crypter
}

func NewAccessToken(cfg config.AccessToken) *AccessToken {
	return &AccessToken{
		cfg:     cfg,
		crypter: NewCrypter(cfg.SignedKey),
	}
}
func (a *AccessToken) generateNewAccessToken(refreshToken string, userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		Subject:   client.Id.String(),
	})

	signedString := refreshToken + a.cfg.SignedKey
	tokenString, err := token.SignedString([]byte(signedString))
	if err != nil {
		return "", fmt.Errorf("AuthUseCases - GenerateToken - token.SignedString: %w", err)
	}

	return "", nil
}
