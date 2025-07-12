package tokens

import (
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/jwt"
	"github.com/go-faster/errors"
)

type refreshToken struct {
	jwtToken jwt.JwtTokenGenerate
	cfg      config.RefreshToken
}

func NewRefreshToken(jwtToken jwt.JwtTokenGenerate, cfg config.RefreshToken) RefreshTokenStore {
	return &refreshToken{
		cfg:      cfg,
		jwtToken: jwtToken,
	}
}

func (a *refreshToken) New() (string, error) {
	return a.jwtToken.GenerateToken(a.cfg.SignedKey, a.cfg.TokenLifetime, "")
}

func (a *refreshToken) Parse(refreshToken string) error {
	if _, err := a.jwtToken.ParseToken(refreshToken, a.cfg.SignedKey); err != nil {
		return errors.Wrap(err, "parseRefreshToken error")
	}
	return nil
}
