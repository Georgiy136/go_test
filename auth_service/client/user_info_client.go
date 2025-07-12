package client

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strconv"
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

func (n *UserInfoClient) GetUserInfo(ctx context.Context, userID int) (*models.User, error) {
	request, err := http.NewRequest(http.MethodGet, n.cfg.Url+"?"+strconv.Itoa(userID), nil)
	if err != nil {
		return nil, fmt.Errorf(" http.NewRequest req error: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(" n.httpClient.Do req error: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var user *models.User
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf(" io.ReadAll req error: %v", err)
		}
		if err = jsoniter.Unmarshal(response, user); err != nil {
			return nil, fmt.Errorf(" jsoniter.Unmarshal req error: %v", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("http status code %d", resp.StatusCode)
}
