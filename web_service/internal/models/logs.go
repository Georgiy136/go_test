package models

import (
	"time"
)

type Log struct {
	Dt           time.Time `json:"dt"`
	Api          string    `json:"api"`
	ServiceName  string    `json:"service_name"`
	Request      string    `json:"request"`
	Response     string    `json:"response"`
	ResponseCode int       `json:"response_code"`
}
