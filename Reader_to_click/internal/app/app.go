package app

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	clickhouse_service "github.com/Georgiy136/go_test/Reader_to_click/internal/clickhouse"
	nats_service "github.com/Georgiy136/go_test/Reader_to_click/internal/nats"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/service"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
	nats_conn "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	click, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Errorf("app - Run - clickhouse.New: %v", err)
	}

	natsConn, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - nats.New: %v", err)
	}

	natsService := nats_service.NewNats(natsConn)
	clickhouse_service.NewClickhouse(click)

	// Запуск крона
	logger := middleware.NewLogger(natsLogs, cfg.Nats.ChannelName)

	// Use case
	goodsUseCases := usecase.NewGoodsUsecases(goodsRepository, goodsRedis)
}
