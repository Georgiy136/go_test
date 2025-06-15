package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"myapp/config"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Postgres struct {
	Conn *bun.DB
	cfg  config.Postgres
}

func New(cfg config.Postgres) (*Postgres, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithDatabase(cfg.Dbname),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
		pgdriver.WithInsecure(true),
	)

	sqlDB := sql.OpenDB(pgconn)
	db := bun.NewDB(sqlDB, pgdialect.New())

	if _, err := db.ExecContext(context.Background(), "select 1"); err != nil {
		return nil, err
	}

	logrus.Debugf("соединение с базой данных postgres успешно установлено")

	return &Postgres{
		Conn: db,
		cfg:  cfg,
	}, nil
}

func (db *Postgres) MigrateUpPostgres() {
	instance, err := postgres.WithInstance(db.Conn.DB, &postgres.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	// Создание мигратора
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		db.cfg.Dbname,
		instance,
	)
	if err != nil {
		logrus.Fatalf("failed to create migrate instance: %v", err)
	}

	// Выполняем миграции
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Debugf("no migration to apply")
		} else {
			logrus.Fatalf("failed to apply migrations: %v", err)
		}
	}

	logrus.Debugf("миграции успешно применены")
}

func (db *Postgres) CloseConn() {
	db.Conn.Close()
}
