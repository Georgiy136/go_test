package app

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/reader"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/service"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	// создаём подключения
	clickConn, err := clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Fatalf("app - Run - clickhouse.New: %v", err)
	}

	// инициализация сервисов
	sendLogsToClick := service.NewSendLogsToClick(clickConn)

	// инициализация reader-ов
	readerService := reader.NewReaderService(cfg.Reader)
	readerService.Configure(
		map[string]reader.HandleFunc{
			"reader_to_click": sendLogsToClick,
		},
	)

	// Запуск
	readerService.Start()
}
