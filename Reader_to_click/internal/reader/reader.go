package reader

import (
	"errors"
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	nats_conn "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"time"
)

type Reader struct {
	handle map[string]HandleFunc
	cfg    config.Reader
}

type HandleFunc interface {
	Run(data [][]byte) error
}

func NewReaderService(cfg config.Reader) *Reader {
	return &Reader{
		cfg: cfg,
	}
}

func (r *Reader) Configure(handle map[string]HandleFunc) {
	r.handle = handle
}

func (r *Reader) Start() {
	for name, handler := range r.handle {
		streamConf, ok := r.cfg.StreamConf[name]
		if !ok {
			logrus.Debugf("conf for handler '%s' not found", name)
			continue
		}
		r.Work(name, handler, streamConf)
	}
}

const (
	timeSleepOnError = 120 * time.Second
	timeSleepOnOk    = 60 * time.Second
)

func (r *Reader) Work(handleName string, handleFunc HandleFunc, readerCfg config.ReaderStreamConf) {
	// подключаем к Nats
	natsConn, err := nats_conn.New(r.cfg.NatsUrl)
	if err != nil {
		logrus.Errorf("[%s]: nats_conn.New error: %v", handleName, err)
		time.Sleep(timeSleepOnError)
	}

	err = createConsumerIfNotExist(natsConn.Js, readerCfg.ChannelName, readerCfg.ConsumerName)
	if err != nil {
		logrus.Errorf("[%s]: can not create subscription: %v", handleName, err)
		time.Sleep(timeSleepOnError)
	}

	sub, err := natsConn.Js.PullSubscribe(readerCfg.ChannelName, readerCfg.ConsumerName, nats.ManualAck())
	if err != nil {
		logrus.Errorf("[%s]: can not create subscription: %v", handleName, err)
		time.Sleep(timeSleepOnError)
	}

	for {
		// Получаем данные из шины
		msgs, err := sub.Fetch(readerCfg.BatchSize)
		if err != nil {
			logrus.Errorf("[%s] error getting msgs, err: %v", handleName, err)
			time.Sleep(timeSleepOnError)
			continue
		}

		var data [][]byte
		for _, msg := range msgs {
			data = append(data, msg.Data)
		}

		if err = handleFunc.Run(data); err != nil {
			logrus.Errorf("[%s] error handling msgs, err: %v", handleName, err)
			time.Sleep(timeSleepOnError)
		}

		for i := range msgs {
			if err = msgs[i].Ack(); err != nil {
				logrus.Errorf("[%s] can not ack msgs: %v", handleName, err)
			}
		}

		time.Sleep(timeSleepOnOk)
	}
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
