package service

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthDBStore interface {
	GetUserSignIn(ctx context.Context, userID int, sessionID string) (*models.LoginInfo, error)
	SaveUserLogin(ctx context.Context, data models.LoginInfo) error
}
