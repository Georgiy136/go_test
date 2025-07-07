package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func GetDataFromDBAndUnmarshal[T any](ctx context.Context, pgconn *pgx.Conn, pg *PgSpec) (*T, error) {
	var result T

	query := pg.GetQuery()
	params := pg.GetParameters()

	if err := pgconn.QueryRow(ctx, query, params...).Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
