package nats

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"myapp/config"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	js  nats.JetStreamContext
	nc  *nats.Conn
	cfg config.Nats
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

	if err = CreateConfStreamIfNotExist(js, cfg.ChannelName); err != nil {
		return nil, fmt.Errorf("Ошибка создания стрима: %v", err)
	}

	logrus.Info("подключение к NATS завершено")

	return &Nats{
		js:  js,
		nc:  nc,
		cfg: cfg,
	}, nil
}

func CreateConfStreamIfNotExist(js nats.JetStreamContext, streamName string) error {
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

func (n *Nats) PublishData(data []byte) error {
	if _, err := n.js.Publish(n.cfg.ChannelName, data); err != nil {
		return fmt.Errorf("error publish to stream: %w, name: %s", err, n.cfg.ChannelName)
	}
	return nil
}

func (n *Nats) Close() {
	n.nc.Close()
}
