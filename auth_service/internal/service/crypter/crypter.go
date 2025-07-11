package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type Crypter struct {
	key []byte
}

func NewCrypter(key string) *Crypter {
	return &Crypter{
		key: []byte(key),
	}
}

// Encrypt зашифровать строку.
func (r *Crypter) Encrypt(payload string) ([]byte, error) {
	block, err := aes.NewCipher(r.key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	text := []byte(payload)

	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("io.ReadFull: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], text)

	return cipherText, nil
}

// Decrypt расшифровать строку.
func (r *Crypter) Decrypt(payload []byte) (string, error) {
	block, err := aes.NewCipher(r.key)
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}

	cipherText := payload[:]

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

// EncryptAndEncodeToBase64 зашифровать и закодировать результат в base64.
func (r *Crypter) EncryptAndEncodeToBase64(payload string) (string, error) {
	encodedStr, err := r.Encrypt(payload)
	if err != nil {
		return "", fmt.Errorf("r.Encrypt: %w", err)
	}

	return r.EncodeToBase64(encodedStr), nil
}

// DecodeFromBase64AndDecrypt декодировать base64 и расшифровать результат.
func (r *Crypter) DecodeFromBase64AndDecrypt(payload string) (string, error) {
	decoded, err := r.DecodeFromBase64(payload)
	if err != nil {
		return "", fmt.Errorf("convhelpers.DecodeFromBase64: %w", err)
	}

	decrypted, err := r.Decrypt(decoded)
	if err != nil {
		return "", fmt.Errorf("r.Decrypt: %w", err)
	}

	return decrypted, nil
}

func (r *Crypter) EncodeToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (r *Crypter) DecodeFromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
