package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Dbpool *pgxpool.Pool
	cfg    config.Postgres
}

const connFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"

func NewPostgres(cfg config.Postgres) (*Postgres, error) {
	dbpool, err := pgxpool.Connect(context.Background(), fmt.Sprintf(connFormat, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname, cfg.Sslmode))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if _, err = dbpool.Exec(context.Background(), "select 1"); err != nil {
		return nil, fmt.Errorf("ping database error: %v", err)
	}

	logrus.Infof("соединение с базой данных postgres успешно установлено")

	return &Postgres{
		Dbpool: dbpool,
		cfg:    cfg,
	}, nil
}

func (db *Postgres) MigrateUpPostgres() error {
	connString := fmt.Sprintf(connFormat, db.cfg.Host, db.cfg.Port, db.cfg.User, db.cfg.Password, db.cfg.Dbname, db.cfg.Sslmode)
	sqlConn, err := sql.Open("postgres", connString)
	if err != nil {
		logrus.Fatal(err)
	}
	defer sqlConn.Close()

	instance, err := postgres.WithInstance(sqlConn, &postgres.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	// Создание мигратора
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations_postgres",
		db.cfg.Dbname,
		instance,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrations to postgres instance: %v", err)
	}

	// Выполняем миграции
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Infof("no migrations to postgres to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations to postgres: %v", err)
	}

	logrus.Infof("migrations to postgres applied")
	return nil
}
