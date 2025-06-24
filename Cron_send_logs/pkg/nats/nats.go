package nats

import (
	"fmt"
	"github.com/Georgiy136/go_test/Cron_send_logs/config"
	"github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Js nats.JetStreamContext
	Nc *nats.Conn
}

func New(cfg config.Nats) (*Nats, error) {
	// Подключение к NATS
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к NATS: %v", err)
	}

	// Создание JetStream контекста
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания JetStream контекста: %v", err)
	}

	logrus.Info("соединение с NATS успешно установлено")

	return &Nats{
		Js: js,
		Nc: nc,
	}, nil
}
