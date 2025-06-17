package repository

import (
	"database/sql"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

func GetDataFromDB[T any](rows *sql.Rows) (*T, error) {
	result := struct {
		Data T `json:"data"`
	}{}

	var data jsoniter.RawMessage
	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("error rows Scan data from db: %v", err)
		}
	}

	if err := jsoniter.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error json unmarshal data from db: %v", err)
	}

	return &result.Data, nil
}

func GetSliceDataFromDB[T any](rows *sql.Rows) (*T, error) {

	result := struct {
		Data T `json:"data"`
	}{}

	return &result.Data, nil
}
