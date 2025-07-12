package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtTokenGenerateGolangJwtV5 struct{}

func NewJwtTokenGenerateGolangJwtV5() JwtTokenGenerate {
	return &JwtTokenGenerateGolangJwtV5{}
}

func (j *JwtTokenGenerateGolangJwtV5) GenerateToken(signedKey string, ttl time.Duration, sub string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		Subject:   sub,
	})

	jwtToken, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", fmt.Errorf("GenerateToken - token_generate.SignedString error: %w", err)
	}

	return jwtToken, nil
}

func (j *JwtTokenGenerateGolangJwtV5) ParseToken(token, signedKey string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired): // если срок действия токена истёк, возвращаем subject и ошибку
			sub, err := j.getTokenSubject(parsedToken)
			if err != nil {
				return "", err
			}
			return sub, err
		default:
			return "", fmt.Errorf("jwt.Parse error: %w", err)
		}
	}

	if parsedToken.Valid {
		return j.getTokenSubject(parsedToken)
	}

	return "", errors.New("token_generate is not valid error")
}

func (j *JwtTokenGenerateGolangJwtV5) getTokenSubject(token *jwt.Token) (string, error) {
	payloadString, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("getTokenSubject error: %w", err)
	}
	return payloadString, err
}
