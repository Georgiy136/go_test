package service

import (
	"github.com/Georgiy136/go_test/auth_service/config"
	"strconv"
)

type RefreshToken struct {
	cfg     config.RefreshToken
	crypter *crypter
}

func NewRefreshToken(cfg config.RefreshToken) *RefreshToken {
	return &RefreshToken{
		cfg:     cfg,
		crypter: NewCrypter(cfg.SignedKey),
	}
}

func (r *RefreshToken) generateNewRefreshToken(userID int) (string, error) {
	refrashToken, err := r.crypter.Encrypt(strconv.Itoa(userID))
	if err != nil {
		return "", err
	}

	return string(refrashToken), nil
}
