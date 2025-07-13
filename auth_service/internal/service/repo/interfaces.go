package repo

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthDBStore interface {
	GetUserSession(ctx context.Context, userID int, sessionID string) (*models.LoginInfo, error)
	SaveUserSession(ctx context.Context, data models.LoginInfo) error
	DeleteUserSession(ctx context.Context, userID int, sessionID string) error
}
