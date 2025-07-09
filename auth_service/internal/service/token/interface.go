package token

import "github.com/Georgiy136/go_test/auth_service/internal/models"

type IssueTokensStore interface {
	GenerateTokensPair(refreshTokenID, userID int) (*models.AuthTokens, error)
}
