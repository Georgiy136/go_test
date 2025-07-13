package models

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenPayload struct {
	UserID    int    `json:"user_id"`
	SessionID string `json:"session_id"`
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

type DataFromRequestGetUser struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DataFromRequestLogout struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginInfo struct {
	UserID    int    `json:"user_id"`
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
	UserAgent string `json:"user_agent"`
	IpAddress string `json:"ip_address"`
}
