package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Goods struct {
	bun.BaseModel `bun:"table:goods"`

	Id          int       `json:"id" bun:"id"`
	ProjectID   int       `json:"project_id" bun:"project_id"`
	Name        string    `json:"name" bun:"name"`
	Description string    `json:"description" bun:"description"`
	Priority    int       `json:"priority" bun:"priority"`
	Removed     bool      `json:"removed" bun:"removed"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at"`
}

type DataFromRequestGoodsAdd struct {
	ProjectID   int
	Name        string
	Description *string
	Priority    *int
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
	ID        int
	ProjectID int
	Limit     int
	Offset    int
}
