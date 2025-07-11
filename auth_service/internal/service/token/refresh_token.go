package token

import (
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"github.com/go-faster/errors"
)

type RefreshToken struct {
	jwtToken jwt.JwtTokenGenerate
	cfg      config.RefreshToken
}

func NewRefreshToken(jwtToken jwt.JwtTokenGenerate, cfg config.RefreshToken) *RefreshToken {
	return &RefreshToken{
		cfg:      cfg,
		jwtToken: jwtToken,
	}
}

func (a *RefreshToken) New() (string, error) {
	return a.jwtToken.GenerateToken(a.cfg.SignedKey, a.cfg.TokenLifetime, nil)
}

func (a *RefreshToken) Parse(refreshToken string) error {
	if _, err := a.jwtToken.ParseToken(refreshToken, a.cfg.SignedKey); err != nil {
		return errors.Wrap(err, "parseRefreshToken error")
	}
	return nil
}
