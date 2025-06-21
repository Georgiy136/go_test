package postgres

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"myapp/internal/errors/common"
)

func getDataFromDB[T any](ctx context.Context, pgconn *pgx.Conn, query string, args ...any) (*T, error) {
	result := struct {
		Data T `json:"data"`
	}{}

	if err := pgconn.QueryRow(ctx, query, args...).Scan(&result); err != nil {
		return nil, ParseProcedureError(err)
	}

	return &result.Data, nil
}

const defaultExceptionErrorCode = "P0001" // Код "P0001" присваивается в случае намеренного возврата ошибки из базы с помощью EXCEPTION (ожидаемая ошибка бизнес логики)

// ParseProcedureError - парсит ошибку из базы данных, которая была инициирована вызовом RAISE EXCEPTION
func ParseProcedureError(procedureErr error) error {
	var pgErr *pgconn.PgError
	if errors.As(procedureErr, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			return &common.BusinessError{Message: pgErr.Error()}
		}
		return procedureErr
	}
	return nil
}
