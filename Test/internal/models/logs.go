package models

import "time"

type Log struct {
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"eventTime"`
}
