package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"myapp/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Pgconn *pgx.Conn
	cfg    config.Postgres
}

const connFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"

/*
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

	logrus.Infof("соединение с базой данных postgres успешно установлено")

	return &Postgres{
		Conn: db,
		cfg:  cfg,
	}, nil
}*/

func NewPostgres(cfg config.Postgres) (*Postgres, error) {
	pgconn, err := pgx.Connect(context.Background(), fmt.Sprintf(connFormat, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname, cfg.Sslmode))
	if err != nil {
		logrus.Fatalf("unable to connect to database: %v\n", err)
	}

	if _, err := pgconn.Exec(context.Background(), "select 1"); err != nil {
		return nil, err
	}

	logrus.Infof("соединение с базой данных postgres успешно установлено")

	return &Postgres{
		Pgconn: pgconn,
		cfg:    cfg,
	}, nil
}

func (db *Postgres) MigrateUpPostgres() error {
	connString := fmt.Sprintf(connFormat, db.cfg.Host, db.cfg.Port, db.cfg.User, db.cfg.Password, db.cfg.Dbname, db.cfg.Sslmode)
	sqlConn, err := sql.Open("postgres", connString)
	if err != nil {
		logrus.Fatal(err)
	}

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

func (db *Postgres) CloseConn() {
	db.Pgconn.Close(context.Background())
}
