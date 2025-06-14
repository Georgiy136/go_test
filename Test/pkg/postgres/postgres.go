package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
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

func migrateUpPostgres(pg *bun.DB) {
	// Накатываем миграции
	driver, err := oci_driver.OracleWithInstance(ctx, ora, &oci_driver.OracleDriverConfig{
		MigrationTable:     "migrations",
		MigrationMutexName: "migrations_mutex",
		NoLock:             true,
		StatementTimeout:   0,
	})

	if err != nil {
		logrus.Fatalf("не удалось создать драйвер для миграций Oracle")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"oracle",
		driver,
	)

	if err != nil {
		logrus.Fatal("не удалось инициализировать систему миграций Oracle")
	}

	err = m.Up()
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatal("ошибка при применении миграций Oracle")
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logrus.Debugf("нет новых миграций для применения")
		return
	}

	logrus.Debugf("миграции успешно применены")
}
