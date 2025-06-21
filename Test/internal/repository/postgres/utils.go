package postgres

import (
	"database/sql"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("rows Scan data from db: %v", string(data))

	if data == nil {
		return nil, nil
	}

	if err := jsoniter.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error json unmarshal data from db: %v", err)
	}
	return &result.Data, nil
}
