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

func NewAuthRepo(pg *postgres.Postgres) service.AuthDBStore {
	return &AuthRepo{
		pgconn: pg.Pgconn,
	}
}

func (db *AuthRepo) SaveUserLogin(ctx context.Context, data models.LoginInfo) error {

	return nil
}

func (db *AuthRepo) GetRefreshTokenID(ctx context.Context) (int, error) {

	return 0, nil
}

func (db *AuthRepo) GetSignInByRefreshTokenID(ctx context.Context, refreshTokenID int) (*models.LoginInfo, error) {

	return &models.LoginInfo{
		UserID:         1,
		RefreshTokenID: refreshTokenID,
		UserAgent:      "IntelliJ HTTP Client/GoLand 2025.1.3",
		IpAddress:      "127.0.0.1",
	}, nil
}
