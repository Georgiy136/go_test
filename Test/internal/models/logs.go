package models

import "time"

type Log struct {
	LogID       int       `json:"id"`
	GoodsID     int       `json:"goods_id"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"eventTime"`
}
