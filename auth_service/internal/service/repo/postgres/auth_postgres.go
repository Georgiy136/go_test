package postgres

import (
	"context"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	"github.com/Georgiy136/go_test/auth_service/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type AuthRepo struct {
	pgconn *pgx.Conn
}

func NewAuthRepo(pg *postgres.Postgres) service.AuthStore {
	return &AuthRepo{
		pgconn: pg.Pgconn,
	}
}
func (db *AuthRepo) GetUser(ctx context.Context, userID int) (*models.User, error) {

	return &models.User{UserID: 1}, nil
}

func (db *AuthRepo) SaveUserLogin(ctx context.Context, data models.LoginInfo) error {

	return nil
}

func (db *AuthRepo) GetRefreshTokenID(ctx context.Context) (int, error) {

	return 0, nil
}

func (db *AuthRepo) GetSignInByRefreshTokenID(ctx context.Context, refreshTokenID int) (*models.LoginInfo, error) {

	return nil, nil
}
