package nats

import (
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/sirupsen/logrus"
	"time"

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

	subject := cfg.ChannelName
	consumerName := cfg.ConsumerName

	sub, err := js.PullSubscribe(subject, consumerName)
	if err != nil {
		return nil, fmt.Errorf("[%s]: can not create subscription: %v", subject, err)
	}

	msgs, err := sub.Fetch(2, nil)
	if err != nil {
		logrus.Errorf("[%s]: error getting msgs on subject: %s, err: %v", subject, err)
		time.Sleep(time.Second)
		continue
	}

	return &Nats{
		Js: js,
		Nc: nc,
	}, nil
}

// createConsumerIfNotExist создает консьюмера на натс-сервере
func createConsumerIfNotExist(js nats.JetStreamContext, streamName, consumerName, subject string) error {
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
		Durable:       consumerName,             // durable имя консьюмера
		Name:          consumerName,             // имя консьюмера (должно быть равно durable)
		AckPolicy:     nats.AckExplicitPolicy,   // требуется подтверждение на все сообщения
		MaxDeliver:    -1,                       // отключаем ограничение на количество перепосылок сообщения
		BackOff:       nil,                      // тоже что и AckWait только позволяет задать промежутки между повтором
		FilterSubject: subject,                  // subject, на который подписывается консьюмер
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
