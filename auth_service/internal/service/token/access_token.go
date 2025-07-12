package token

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token/jwt"
	"github.com/go-faster/errors"
	jsoniter "github.com/json-iterator/go"
)

type AccessToken struct {
	jwtToken jwt.JwtTokenGenerate
	cfg      config.AccessToken
}

func NewAccessToken(jwtToken jwt.JwtTokenGenerate, cfg config.AccessToken) *AccessToken {
	return &AccessToken{
		cfg:      cfg,
		jwtToken: jwtToken,
	}
}

func (a *AccessToken) New(refreshToken string, accessTokenPayload models.AccessTokenPayload) (string, error) {
	return a.jwtToken.GenerateToken(a.getSignedString(refreshToken), a.cfg.TokenLifetime, accessTokenPayload)
}

func (a *AccessToken) Parse(tokens models.AuthTokens) (*models.AccessTokenPayload, error) {
	payloadString, err := a.jwtToken.ParseToken(tokens.AccessToken, a.getSignedString(tokens.RefreshToken))
	if err != nil {
		switch {
		case errors.Is(err, jwt.TokenIsExpiredError):
			payload, err := a.getAccessTokenPayload(payloadString)
			if err != nil {
				return payload, fmt.Errorf("ParseToken - a.getAccessTokenPayload error: %w", err)
			}
			return payload, err
		}
		return nil, errors.Wrap(err, "ParseToken error")
	}

	payload, err := a.getAccessTokenPayload(payloadString)
	if err != nil {
		return nil, fmt.Errorf("ParseToken - a.getAccessTokenPayload error: %w", err)
	}
	return payload, nil
}

func (a *AccessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}

func (a *AccessToken) getAccessTokenPayload(payload string) (*models.AccessTokenPayload, error) {
	var payloadData models.AccessTokenPayload
	if err := jsoniter.UnmarshalFromString(payload, &payloadData); err != nil {
		return nil, err
	}
	return &payloadData, nil
}
