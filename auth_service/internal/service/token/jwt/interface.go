package jwt

import "time"

type JwtTokenGenerate interface {
	GenerateToken(signedKey string, ttl time.Duration, payload string) (string, error)
	ParseToken(token, signedKey string) (string, error)
}
