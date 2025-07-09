package models

import "time"

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInfo struct {
	Issuer         string    `json:"issuer"`
	RefreshTokenID int       `json:"refresh_token_id"`
	ExpiredAt      time.Time `json:"expired_at"`
	IssuedAt       time.Time `json:"issued_at"`
}

type AccessTokenInfo struct {
	Issuer    string             `json:"issuer"`
	Payload   AccessTokenPayload `json:"payload"`
	ExpiredAt time.Time          `json:"expired_at"`
	IssuedAt  time.Time          `json:"issued_at"`
}

type AccessTokenPayload struct {
	UserID         int `json:"user_id"`
	RefreshTokenID int `json:"refresh_token_id"`
}

type DataFromRequestGetTokens struct {
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	IpAddress string `json:"ip_address"`
}

type DataFromRequestUpdateTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IpAddress    string `json:"ip_address"`
}

type SaveLoginDataDbRequest struct {
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IpAddress    string `json:"ip_address"`
}
