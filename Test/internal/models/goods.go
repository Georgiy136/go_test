package models

import "time"

type GoodsUpdDBResponse struct {
	Data *Goods `json:"data"`
}

type Goods struct {
	GoodID      int        `json:"good_id"`
	ProjectID   int        `json:"project_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Priority    int        `json:"priority"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type DataFromRequestGoodsAdd struct {
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
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

type GoodsListDBResponse struct {
	Data *GoodsList `json:"data"`
}

type GoodsList struct {
	Meta *struct {
		Total  int `json:"total"`
		Remove int `json:"remove"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta,omitempty"`
	Goods []Goods `json:"goods"`
}

type DataFromRequestReprioritizeGood struct {
	GoodsID   int `json:"good_id"`
	ProjectID int `json:"project_id"`
	Priority  int `json:"priority"`
}
