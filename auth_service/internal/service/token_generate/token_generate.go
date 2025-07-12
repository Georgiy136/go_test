package token_generate

import (
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/tokens"
)

type IssueTokensService struct {
	RefreshToken tokens.RefreshTokenStore
	AccessToken  tokens.AccessTokenStore
	crypter      *crypter.Crypter
}

func NewIssueTokensService(
	refreshToken tokens.RefreshTokenStore,
	accessToken tokens.AccessTokenStore,
	crypter *crypter.Crypter,
) *IssueTokensService {
	return &IssueTokensService{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		crypter:      crypter,
	}
}
