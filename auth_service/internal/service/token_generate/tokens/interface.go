package tokens

import "github.com/Georgiy136/go_test/auth_service/internal/models"

type (
	AccessTokenStore interface {
		New(refreshToken string, accessTokenPayload models.AccessTokenPayload) (string, error)
		Parse(tokens models.AuthTokens) (*models.AccessTokenPayload, error)
	}

	RefreshTokenStore interface {
		New() (string, error)
		Parse(refreshToken string) error
	}
)
