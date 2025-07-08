package token

import "github.com/Georgiy136/go_test/auth_service/internal/models"

type IssueTokensStore interface {
	GenerateTokensPair(userID int) (*models.AuthTokens, error)
}
