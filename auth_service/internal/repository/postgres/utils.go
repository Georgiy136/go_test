package postgres

import (
	"context"
	"github.com/Georgiy136/go_test/web_service/internal/errors/common"
	"github.com/Georgiy136/go_test/web_service/pkg/postgres"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetDataFromDBAndUnmarshal[T any](ctx context.Context, pgconn *pgx.Conn, pg *postgres.PgSpec) (*T, error) {
	var result T

	query := pg.GetQuery()
	params := pg.GetParameters()

	if err := pgconn.QueryRow(ctx, query, params...).Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, handleDBError(err)
	}

	return &result, nil
}

const defaultExceptionErrorCode = "P0001" // Код "P0001" присваивается в случае намеренного возврата ошибки из базы с помощью EXCEPTION (ожидаемая ошибка бизнес логики)

func handleDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			return &common.CustomError{Description: pgErr.Message, Err: &common.ServiceUnprocessableEntity}
		}
		return err
	}
	return err
}
