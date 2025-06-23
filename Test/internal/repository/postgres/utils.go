package postgres

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"myapp/internal/errors/common"
)

func GetDataFromDB[T any](ctx context.Context, pgconn *pgx.Conn, query string, args ...any) (*T, error) {
	var result T

	if err := pgconn.QueryRow(ctx, query, args...).Scan(&result); err != nil {
		return nil, fmt.Errorf("db error: %w", HandleDBError(err))
	}

	return &result, nil
}

const defaultExceptionErrorCode = "P0001" // Код "P0001" присваивается в случае намеренного возврата ошибки из базы с помощью EXCEPTION (ожидаемая ошибка бизнес логики)

func HandleDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			return &common.CustomError{Message: pgErr.Message, Err: &common.ServiceUnprocessableEntity}
		}
		return err
	}
	return err
}
