package client

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"net/http"
)

type UserInfoClient struct {
	httpClient *http.Client
	cfg        config.UserInfoClient
}

func NewUserInfoClient(cfg config.UserInfoClient) *UserInfoClient {
	return &UserInfoClient{
		httpClient: http.DefaultClient,
		cfg:        cfg,
	}
}

// mock функция
func (n *UserInfoClient) GetUserInfo(ctx context.Context, userID int) (*models.User, error) {

	return &models.User{UserID: userID}, nil
}
