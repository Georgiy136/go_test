package postgres

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"myapp/internal/errors/common"
	"myapp/pkg/postgres"
)

func GetDataFromDB[T any](ctx context.Context, pgconn *pgx.Conn, pg *postgres.PgSpec) (*T, error) {
	var result T

	query := pg.GetQuery()
	params := pg.GetParameters()

	if err := pgconn.QueryRow(ctx, query, params...).Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("db error: %w", HandleDBError(err))
	}

	return &result, nil
}

const defaultExceptionErrorCode = "P0001" // Код "P0001" присваивается в случае намеренного возврата ошибки из базы с помощью EXCEPTION (ожидаемая ошибка бизнес логики)

func HandleDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			return &common.CustomError{Description: pgErr.Message, Err: &common.ServiceUnprocessableEntity}
		}
		return err
	}
	return err
}
