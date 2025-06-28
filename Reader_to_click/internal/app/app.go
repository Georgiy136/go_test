package app

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	clickhouse_service "github.com/Georgiy136/go_test/Reader_to_click/internal/clickhouse"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/cron"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
	nats_conn "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	clickConn, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Errorf("app - Run - clickhouse.New: %v", err)
	}

	natsConn, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - cron.New: %v", err)
	}
	click := clickhouse_service.NewClickhouse(clickConn)

	// инициализация сервисов
	cron := cron.NewCron(cfg.Cron, natsConn, click)

	// Запуск сервисов
	cron.Start()
}
