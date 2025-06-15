package clickhouse

import (
	"context"
	"crypto/tls"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/go-clickhouse/ch"
	"myapp/config"
	"time"
)

type Clickhouse struct {
	Conn *ch.DB
	cfg  config.Clickhouse
}

func New(cfg config.Clickhouse) (*Clickhouse, error) {
	db := ch.Connect(
		ch.WithAddr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		ch.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		ch.WithUser(cfg.User),
		ch.WithPassword(cfg.Password),
		ch.WithDatabase(cfg.Dbname),
		ch.WithTimeout(5*time.Second),
		ch.WithDialTimeout(5*time.Second),
		ch.WithReadTimeout(5*time.Second),
		ch.WithWriteTimeout(5*time.Second),
		ch.WithQuerySettings(map[string]interface{}{
			"prefer_column_name_to_alias": 1,
		}),
	)

	//sqlDB := sql.OpenDB(db)

	if _, err := db.ExecContext(context.Background(), "select 1"); err != nil {
		return nil, err
	}

	logrus.Infof("соединение с базой данных clickhouse успешно установлено")

	return &Clickhouse{
		Conn: db,
		cfg:  cfg,
	}, nil
}

/*
func (db *Clickhouse) MigrateUpClickhouse() error {
	instance, err := clickhouse.WithInstance(db.Conn, &clickhouse.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	// Создание мигратора
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations_clickhouse",
		db.cfg.Dbname,
		instance,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrations_clickhouse instance: %v", err)
	}

	// Выполняем миграции
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Infof("no migrations to clickhouse to apply")
			return nil
		}
		return fmt.Errorf("failed to apply migrations to clickhouse: %v", err)
	}

	logrus.Infof("migrations to clickhouse applied")
	return nil
}*/

func (db *Clickhouse) CloseConn() {
	db.Conn.Close()
}
