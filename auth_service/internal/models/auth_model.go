package models

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

type LoginInfo struct {
	UserID         int    `json:"user_id"`
	RefreshTokenID int    `json:"refresh_token_id"`
	RefreshToken   string `json:"refresh_token"`
	UserAgent      string `json:"user_agent"`
	IpAddress      string `json:"ip_address"`
}
