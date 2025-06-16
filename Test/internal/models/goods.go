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
