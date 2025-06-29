package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Nats struct {
	Js nats.JetStreamContext
	Nc *nats.Conn
}

func New(natsURL string) (*Nats, error) {
	// Подключение к NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS: %v", err)
	}

	// Создание JetStream контекста
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания JetStream контекста: %v", err)
	}

	logrus.Info("соединение с NATS успешно установлено")

	return &Nats{
		Js: js,
		Nc: nc,
	}, nil
}
