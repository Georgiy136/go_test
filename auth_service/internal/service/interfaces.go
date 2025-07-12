package service

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthDBStore interface {
	GetRefreshTokenID(ctx context.Context) (int, error)
	GetSignInByRefreshTokenID(ctx context.Context, refreshTokenID int) (*models.LoginInfo, error)
	SaveUserLogin(ctx context.Context, data models.LoginInfo) error
}
