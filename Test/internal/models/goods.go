package models

import "time"

type Goods struct {
	GoodID      int       `json:"good_id"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	DeletedAt   time.Time `json:"deleted_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type DataFromRequestGoodsAdd struct {
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

type DataFromRequestGoodsUpdate struct {
	ID          int
	ProjectID   int
	Name        string
	Description *string
}

type DataFromRequestGoodsDelete struct {
	ID        int
	ProjectID int
}
type DataFromRequestGoodsList struct {
	GoodsID   int `json:"goods_id"`
	ProjectID int `json:"project_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}

type GoodsListDBResponse struct {
	Meta struct {
		Total  int `json:"total"`
		Remove int `json:"remove"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Goods []Goods `json:"goods"`
}

type DataFromRequestReprioritizeGood struct {
	ID          int
	NewPriority int
}
