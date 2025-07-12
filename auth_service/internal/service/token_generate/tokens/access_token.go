package tokens

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/jwt"
	"github.com/go-faster/errors"
	jsoniter "github.com/json-iterator/go"
)

type accessToken struct {
	jwtToken jwt.JwtTokenGenerate
	crypter  *crypter.Crypter
	cfg      config.AccessToken
}

func NewAccessToken(jwtToken jwt.JwtTokenGenerate, crypter *crypter.Crypter, cfg config.AccessToken) AccessTokenStore {
	return &accessToken{
		jwtToken: jwtToken,
		crypter:  crypter,
		cfg:      cfg,
	}
}

func (a *accessToken) New(refreshToken string, accessTokenPayload models.AccessTokenPayload) (string, error) {
	payloadString, err := a.genAccessTokenPayload(accessTokenPayload)
	if err != nil {
		return "", fmt.Errorf("genAccessTokenPayload error: %w", err)
	}

	return a.jwtToken.GenerateToken(a.getSignedString(refreshToken), a.cfg.TokenLifetime, payloadString)
}

func (a *accessToken) Parse(tokens models.AuthTokens) (*models.AccessTokenPayload, error) {
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

func (a *accessToken) getSignedString(refreshToken string) string {
	return refreshToken + a.cfg.SignedKey
}

func (a *accessToken) genAccessTokenPayload(accessTokenPayload models.AccessTokenPayload) (string, error) {
	payloadString, err := jsoniter.MarshalToString(accessTokenPayload)
	if err != nil {
		return "", fmt.Errorf("ParseToken - a.getAccessTokenPayload error: %w", err)
	}
	payloadEncrypted, err := a.crypter.EncryptAndEncodeToBase64(payloadString)
	if err != nil {
		return "", fmt.Errorf("a.crypter.Encrypt accessToken error: %w", err)
	}
	return payloadEncrypted, nil
}

func (a *accessToken) getAccessTokenPayload(payload string) (*models.AccessTokenPayload, error) {
	payloadDecoded, err := a.crypter.DecodeFromBase64AndDecrypt(payload)
	if err != nil {
		return nil, fmt.Errorf("UpdateTokens - DecodeFromBase64AndDecrypt refreshToken error: %w", err)
	}
	var payloadData models.AccessTokenPayload
	if err = jsoniter.UnmarshalFromString(payloadDecoded, &payloadData); err != nil {
		return nil, err
	}
	return &payloadData, nil
}
