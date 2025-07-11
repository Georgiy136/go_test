package token

import (
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"github.com/go-faster/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type RefreshToken struct {
	jwtToken      jwt.JwtTokenGenerate
	cfg           config.RefreshToken
	tokenLifetime time.Duration
}

func NewRefreshToken(jwtToken jwt.JwtTokenGenerate, cfg config.RefreshToken) *RefreshToken {
	tokenLifetime, err := time.ParseDuration(cfg.TokenLifetime)
	if err != nil {
		logrus.Fatalf("NewRefreshToken: tokenLifetime ParseDuration err: %v", err)
	}

	return &RefreshToken{
		cfg:           cfg,
		tokenLifetime: tokenLifetime,
		jwtToken:      jwtToken,
	}
}
func (a *RefreshToken) generateNewRefreshToken() (string, error) {
	return a.jwtToken.GenerateToken(a.getSignedString(), a.tokenLifetime, nil)
}

func (a *RefreshToken) parseRefreshToken(refreshToken string) error {
	if _, err := a.jwtToken.ParseToken(refreshToken, a.getSignedString()); err != nil {
		return errors.Wrap(err, "parseRefreshToken error")
	}
	return nil
}

func (a *RefreshToken) getSignedString() string {
	return a.cfg.SignedKey
}
