package app

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/reader"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/service"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
	nats_conn "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	// создаём подключения
	clickConn, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Fatalf("app - Run - clickhouse.New: %v", err)
	}
	natsConn, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Fatalf("app - Run - nats_conn.New: %v", err)
	}

	// инициализация сервисов
	service := service.NewService(clickConn)

	reader := reader.NewReader(cfg.Cron, natsConn, service)

	// Запуск
	reader.Start()
}
