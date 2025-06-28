package nats

import (
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Nats struct {
	Js  nats.JetStreamContext
	Nc  *nats.Conn
	Sub *nats.Subscription
}

func New(cfg config.Nats) (*Nats, error) {
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

	logrus.Info("соединение с NATS успешно установлено")

	err = createConsumerIfNotExist(js, cfg.ChannelName, cfg.ConsumerName)
	if err != nil {
		return nil, fmt.Errorf("[%s]: can not create subscription: %v", cfg.ChannelName, err)
	}

	sub, err := js.PullSubscribe(cfg.ChannelName, cfg.ConsumerName, nats.ManualAck())
	if err != nil {
		return nil, fmt.Errorf("[%s]: can not create subscription: %v", cfg.ChannelName, err)
	}

	return &Nats{
		Js:  js,
		Nc:  nc,
		Sub: sub,
	}, nil
}

func createConsumerIfNotExist(js nats.JetStreamContext, streamName, consumerName string) error {
	_, err := js.ConsumerInfo(streamName, consumerName)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, nats.ErrConsumerNotFound):
		break
	default:
		return fmt.Errorf("can not get consumer info: %w", err)
	}

	consumerConfig := nats.ConsumerConfig{
		Durable:    consumerName,           // durable имя консьюмера
		Name:       consumerName,           // имя консьюмера (должно быть равно durable)
		AckPolicy:  nats.AckExplicitPolicy, // требуется подтверждение на все сообщения
		MaxDeliver: -1,                     // отключаем ограничение на количество перепосылок сообщения
		BackOff:    nil,                    // тоже что и AckWait только позволяет задать промежутки между повтором
		// FilterSubject: subject,               // subject, на который подписывается консьюмер
		ReplayPolicy:  nats.ReplayInstantPolicy, // прокачка как можно быстрее, ReplayOriginalPolicy - в оригинальном темпе
		MaxAckPending: 0,                        // максимальное значение сообщений ожидающих Ack
		Replicas:      1,                        // количество реплик консьюмера
		MemoryStorage: false,                    // заставляет сервер хранить состояние косьюмера в ОЗУ
	}

	_, err = js.AddConsumer(streamName, &consumerConfig)
	if err != nil {
		return fmt.Errorf("error creating consumer, stream: %s, consumer: %s, err: %w", streamName, consumerName, err)
	}

	return nil
}
