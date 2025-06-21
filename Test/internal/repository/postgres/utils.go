package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

const (
	// Код "P0001" присваивается в случае намеренного возврата ошибки из базы
	// с помощью EXCEPTION (ожидаемая ошибка бизнес логики)
	defaultExceptionErrorCode = "P0001"
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

// ParseProcedureError - парсит ошибку из базы данных, которая была инициирована вызовом RAISE EXCEPTION
func ParseProcedureError(procedureErr error) error {
	var pgErr *pgconn.PgError
	if errors.As(procedureErr, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			restDataCustomErr, parseErr := parseExceptionPgErrorMessage(pgErr)
			if parseErr != nil {
				return setRestDataErrorParseProcedureError(parseErr)
			}

			return rest_data.NewCustomError(
				errors.New(restDataCustomErr.Detail),
				restDataCustomErr.ErrorKey,
				restDataCustomErr.Message,
				restDataCustomErr.Detail,
				http.StatusUnprocessableEntity,
			)
		}
		return setRestDataReadResponseFromDb(errors.New(pgErr.Error()))
	}

	return setRestDataErrorParseProcedureError(procedureErr)
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
