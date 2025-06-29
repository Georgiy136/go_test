package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	click "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"net"
)

type Clickhouse struct {
	Conn *sql.DB
	Cfg  config.Clickhouse
}

func New(cfg config.Clickhouse) (*Clickhouse, error) {
	conn := click.OpenDB(&click.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: click.Auth{
			Database: cfg.Dbname,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
	})
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("clickhouse Ping failed: %w", err)
	}

	logrus.Infof("соединение с базой данных clickhouse успешно установлено")

	return &Clickhouse{
		Conn: conn,
		Cfg:  cfg,
	}, nil
}

func (db *Clickhouse) MigrateUpClickhouse() error {
	instance, err := clickhouse.WithInstance(db.Conn, &clickhouse.Config{})
	if err != nil {
		return fmt.Errorf("MigrateUpClickhouse - clickhouse.WithInstance failed: %w", err)
	}

	// Создание мигратора
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations_clickhouse",
		db.Cfg.Dbname,
		instance,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrations_clickhouse instance: %w", err)
	}

	// Выполняем миграции
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Infof("no migrations to clickhouse to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations to clickhouse: %w", err)
	}

	logrus.Infof("migrations to clickhouse applied")

	return nil
}
