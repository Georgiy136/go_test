package client

import (
	"bytes"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type NotificationClient struct {
	cfg        config.NotificationClient
	httpClient *http.Client
}

func NewNotificationClient(cfg config.NotificationClient) *NotificationClient {
	return &NotificationClient{
		httpClient: http.DefaultClient,
		cfg:        cfg,
	}
}

func (n *NotificationClient) SendNewSignInNotification(userID int, msg string) error {
	req := struct {
		UserId  int    `json:"user_id"`
		Message string `json:"message"`
	}{
		userID,
		msg,
	}

	body, err := jsoniter.Marshal(req)
	if err != nil {
		return fmt.Errorf(" jsoniter.Marshal notification body req error: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, n.cfg.Url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf(" http.NewRequest req error: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(" n.httpClient.Do req error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code %d", resp.StatusCode)
	}

	return nil
}
