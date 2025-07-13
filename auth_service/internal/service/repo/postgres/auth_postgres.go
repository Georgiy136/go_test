package postgres

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service/app_errors"
	"github.com/Georgiy136/go_test/auth_service/internal/service/repo"
	"github.com/Georgiy136/go_test/auth_service/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthRepo struct {
	Dbpool *pgxpool.Pool
}

func NewAuthRepo(pg *postgres.Postgres) repo.AuthDBStore {
	return &AuthRepo{
		Dbpool: pg.Dbpool,
	}
}

func (db *AuthRepo) SaveUserSession(ctx context.Context, data models.LoginInfo) error {
	conn, err := db.Dbpool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()

	query := `INSERT INTO sessions.user_login (user_id, 
                                 			   session_id,
                                 			   hash_token,
                                 			   user_agent, 
                                 			   ip_address
                                 			   ) 
				values ($1, $2, $3, $4, $5);`

	if _, err = conn.Query(ctx, query, data.UserID, data.SessionID, data.Token, data.UserAgent, data.IpAddress); err != nil {
		return fmt.Errorf("SaveUserSession err: %v", err)
	}
	return nil
}

func (db *AuthRepo) GetUserSession(ctx context.Context, userID int, sessionID string) (*models.LoginInfo, error) {
	conn, err := db.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()

	query := `SELECT user_id,
                     session_id,
                     hash_token,
                     user_agent,
                     ip_address 
			  FROM sessions.user_login
			  WHERE user_id = $1 AND session_id = $2;`

	var result models.LoginInfo
	err = conn.QueryRow(ctx, query, userID, sessionID).Scan(&result.UserID, &result.SessionID, &result.Token, &result.UserAgent, &result.IpAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, app_errors.SessionUserNotFoundError
		}
		return nil, fmt.Errorf("GetUserSession err: %v", err)
	}
	return &result, nil
}

func (db *AuthRepo) DeleteUserSession(ctx context.Context, userID int, sessionID string) error {
	conn, err := db.Dbpool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()

	query := `DELETE FROM FROM sessions.user_login WHERE user_id = $1 AND session_id = $2;`

	if _, err = conn.Query(ctx, query, userID, sessionID); err != nil {
		return fmt.Errorf("DeleteUserSession err: %v", err)
	}
	return nil
}
