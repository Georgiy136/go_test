package service

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
)

type AuthStrore interface {
	GetUser(ctx context.Context, userID int) (*models.User, error)
	SaveUserLogin(ctx context.Context, data models.SaveLoginDataDbRequest) error
}
