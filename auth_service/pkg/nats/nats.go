package nats

import (
	"context"
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/web_service/config"
	"github.com/Georgiy136/go_test/web_service/pkg/jaegerotel"
	"github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Js  nats.JetStreamContext
	Nc  *nats.Conn
	Cfg config.Nats
}

func New(tctx context.Context, cfg config.Nats) (*Nats, error) {
	_, span := jaegerotel.StartSpan(tctx, "Nats - connect")
	defer span.End()

	// Подключение к NATS
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS: %v", err)
	}

	// Создание JetStream контекста
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания JetStream контекста: %v", err)
	}

	if err = CreateStreamIfNotExist(js, cfg.ChannelName); err != nil {
		return nil, fmt.Errorf("ошибка создания стрима: %v", err)
	}

	logrus.Info("соединение с NATS успешно установлено")

	return &Nats{
		Js:  js,
		Nc:  nc,
		Cfg: cfg,
	}, nil
}

func CreateStreamIfNotExist(js nats.JetStreamContext, streamName string) error {
	_, err := js.StreamInfo(streamName)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, nats.ErrStreamNotFound):
		break
	default:
		return fmt.Errorf("can not get stream info: %w", err)
	}

	cfg := nats.StreamConfig{
		Name:              streamName,
		Storage:           nats.FileStorage,
		Retention:         nats.LimitsPolicy,
		MaxConsumers:      -1,
		MaxMsgs:           -1,
		MaxBytes:          -1,
		MaxMsgsPerSubject: -1,
		Discard:           nats.DiscardOld,
	}

	_, err = js.AddStream(&cfg)
	if err != nil {
		return fmt.Errorf("can not create stream: %w, name: %s", err, streamName)
	}

	return nil
}
