package token

import "github.com/Georgiy136/go_test/auth_service/internal/models"

type IssueTokensStore interface {
	GenerateTokensPair(data models.TokenPayload) (*models.AuthTokens, error)
}
