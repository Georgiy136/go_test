package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtTokenGenerateGolangJwtV5 struct{}

func NewJwtTokenGenerateGolangJwtV5() JwtTokenGenerate {
	return &JwtTokenGenerateGolangJwtV5{}
}

func (j *JwtTokenGenerateGolangJwtV5) GenerateToken(signedKey string, ttl time.Duration, payload string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		Subject:   payload,
	})

	jwtToken, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", fmt.Errorf("GenerateToken - token.SignedString error: %w", err)
	}

	return jwtToken, nil
}

func (j *JwtTokenGenerateGolangJwtV5) ParseToken(token, signedKey string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("ParseToken - jwt.Parse error: %w", err)
	}

	if parsedToken.Valid {
		payloadString, err := parsedToken.Claims.GetSubject()
		if err != nil {
			return "", fmt.Errorf("decodeAccessToken claims.GetSubject error: %w", err)
		}

		return payloadString, nil
	}

	return "", fmt.Errorf("ParseToken error, parsedToken not valid: %w", err)
}
