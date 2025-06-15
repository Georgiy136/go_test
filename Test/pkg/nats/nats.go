package nats

import (
	"fmt"
	"github.com/uptrace/bun"
	"log"
	"myapp/config"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Conn *bun.DB
	cfg  config.Nats
}

func New(cfg config.Nats) (*Nats, error) {
	// Подключение к NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer nc.Close()

	// Создание JetStream контекста
	_, err = nc.JetStream()
	if err != nil {
		log.Fatalf("Ошибка создания JetStream контекста: %v", err)
	}

	/*
		// Создание стрима
		streamName := "MY_STREAM"
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"my.subject"},
		})
		if err != nil {
			log.Fatalf("Ошибка создания стрима: %v", err)
		}

		// Публикация сообщения в стрим
		msg := []byte("Hello, JetStream!")
		_, err = js.Publish("my.subject", msg)
		if err != nil {
			log.Fatalf("Ошибка публикации сообщения: %v", err)
		}

		fmt.Println("Сообщение успешно опубликовано в JetStream!")*/

	return &Nats{nil, cfg}, nil
}
