package models

import (
	"time"
)

type Log struct {
	Dt           time.Time `clickhouse:"dt" json:"dt"`
	Api          string    `clickhouse:"api" json:"api"`
	ServiceName  string    `clickhouse:"service_name" json:"service_name"`
	Request      string    `clickhouse:"request" json:"request"`
	Response     string    `clickhouse:"response" json:"response"`
	ResponseCode int       `clickhouse:"response_code" json:"response_code"`
}
