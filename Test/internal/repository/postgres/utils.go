package postgres

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"myapp/pkg/postgres"
)

func getDataFromDB[T any](ctx context.Context, pgconn *pgx.Conn, pg *postgres.PgSpec) (*T, error) {
	var result T

	query := pg.GetQuery()
	params := pg.GetParameters()

	if err := pgconn.QueryRow(ctx, query, params...).Scan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const defaultExceptionErrorCode = "P0001" // Код "P0001" присваивается в случае намеренного возврата ошибки из базы с помощью EXCEPTION (ожидаемая ошибка бизнес логики)

// ProcedureError - парсит ошибку из базы данных, которая была инициирована вызовом RAISE EXCEPTION
func IsProcedureError(err error) (bool, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == defaultExceptionErrorCode {
			return true, pgErr.Error()
		}
		return false, ""
	}
	return false, ""
}
