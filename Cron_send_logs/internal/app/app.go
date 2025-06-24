package app

import (
	"github.com/Georgiy136/go_test/Cron_send_logs/config"
	nats_service "github.com/Georgiy136/go_test/Cron_send_logs/internal/nats"
	"github.com/Georgiy136/go_test/Cron_send_logs/internal/service"
	"github.com/Georgiy136/go_test/Cron_send_logs/pkg/clickhouse"
	nats_conn "github.com/Georgiy136/go_test/Cron_send_logs/pkg/nats"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	// connections
	_, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Errorf("app - Run - clickhouse.New: %v", err)
	}

	nats, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - nats.New: %v", err)
	}

	// очередь Nats для сохранения логов
	natsLogs := nats_service.NewNats(nats)

	// Запуск крона
	logger := middleware.NewLogger(natsLogs, cfg.Nats.ChannelName)

	// Use case
	goodsUseCases := usecase.NewGoodsUsecases(goodsRepository, goodsRedis)
}
