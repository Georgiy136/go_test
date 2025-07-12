package postgres

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	"github.com/Georgiy136/go_test/auth_service/internal/service/app_errors"
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
	query := `INSERT INTO sessions.user_login (user_id, 
                                 			   session_id,
                                 			   hash_refresh_token,
                                 			   user_agent, 
                                 			   ip_address
                                 			   ) 
				values ($1, $2, $3, $4, $5);`

	_, err := db.pgconn.Query(ctx, query, data.UserID, data.SessionID, data.RefreshToken, data.UserAgent, data.IpAddress)
	if err != nil {
		return fmt.Errorf("SaveUserLogin err: %v", err)
	}
	return nil
}

func (db *AuthRepo) GetUserSignIn(ctx context.Context, userID int, sessionID string) (*models.LoginInfo, error) {
	query := `SELECT user_id,
                     session_id,
                     hash_refresh_token,
                     user_agent,
                     ip_address 
			  FROM sessions.user_login
			  WHERE user_id = $1 AND session_id = $2;`

	var result models.LoginInfo
	err := db.pgconn.QueryRow(ctx, query, userID, sessionID).Scan(&result.UserID, &result.SessionID, &result.RefreshToken, &result.UserAgent, &result.IpAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, app_errors.SessionUserNotFoundError
		}
		return nil, fmt.Errorf("GetUserSignIn err: %v", err)
	}
	return &result, nil
}
