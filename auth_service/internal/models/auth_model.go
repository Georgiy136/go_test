package models

import "time"

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	Issuer    string    `json:"issuer"`
	Payload   string    `json:"payload"`
	ExpiredAt time.Time `json:"expired_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

type DataFromRequestGetTokens struct {
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	IpAddress string `json:"ip_address"`
}

type SaveLoginDataDbRequest struct {
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IpAddress    string `json:"ip_address"`
}

type DataFromRequestGoodsUpdate struct {
	GoodID      int     `json:"good_id"`
	ProjectID   int     `json:"project_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type DataFromRequestGoodsDelete struct {
	GoodID    int       `json:"good_id"`
	ProjectID int       `json:"project_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
type DataFromRequestGoodsList struct {
	GoodsID   *int
	ProjectID *int
	Limit     *int
	Offset    *int
}

type DataFromRequestReprioritizeGood struct {
	GoodID    int `json:"good_id"`
	ProjectID int `json:"project_id"`
	Priority  int `json:"priority"`
}
