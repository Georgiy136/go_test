package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func getDataFromDB[T any](ctx context.Context, pgconn *pgx.Conn, query string, args ...any) (*T, error) {
	result := struct {
		Data T `json:"data"`
	}{}

	err := pgconn.QueryRow(ctx, query, args...).Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("getDataFromDB ERROR: %w", err)
	}

	return &result.Data, nil
}

/*
func GetDataFromDB[T any](rows *pgx.Rows) (*T, error) {
	result := struct {
		Data T `json:"data"`
	}{}

	data, err := pgx.CollectOneRow(rows, result)

	//var data jsoniter.RawMessage
	//for rows.Next() {
	//	if err := rows.Scan(&data); err != nil {
	//		return nil, fmt.Errorf("error rows Scan data from db: %v", err)
	//	}
	//}
	logrus.Infof("rows Scan data from db: %v", string(data))

	if data == nil {
		return nil, nil
	}

	if err := jsoniter.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error json unmarshal data from db: %v", err)
	}
	return &result.Data, nil
}*/
