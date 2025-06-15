package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"myapp/config"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

func New(cfg config.Postgres) (*bun.DB, error) {
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

	return db, nil
}

func MigrateUpPostgres(pg *bun.DB, cfg config.Postgres) {
	pg.Conn(context.Background())

	log.Printf("pg.String(): %v", pg.String())

	// Создание мигратора
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf(pg.String()),
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	// Выполняем миграции
	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no migration to apply")
		} else {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	}

	logrus.Debugf("миграции успешно применены")
}
