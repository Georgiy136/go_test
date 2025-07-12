package service

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthDBStore interface {
	GetSignInByTokenID(ctx context.Context, refreshTokenID int) (*models.LoginInfo, error)
	SaveUserLogin(ctx context.Context, data models.LoginInfo) error
	SaveToken(ctx context.Context, refreshToken string) (int, error)
}
