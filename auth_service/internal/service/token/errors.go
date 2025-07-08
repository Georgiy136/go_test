package service

import "errors"

var (
	ErrCipherTextIsTooShort = errors.New("ciphertext is too short")
)
