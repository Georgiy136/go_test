package reader

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	nats_pkg "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
	"time"
)

type Reader struct {
	cfg    config.Cron
	nats   *nats_pkg.Nats
	handle map[string]HandleFunc
}

type HandleFunc interface {
	Run(data [][]byte) error
}

func NewReaderService(cfg config.Cron, nats *nats_pkg.Nats) *Reader {
	return &Reader{
		cfg:  cfg,
		nats: nats,
	}
}

func (r *Reader) Configure(handle map[string]HandleFunc) {
	r.handle = handle
}

func (r *Reader) Start() {
	for name, handler := range r.handle {
		r.Work(name, handler)
	}
}

func (r *Reader) Work(handleName string, handleFunc HandleFunc) {
	const batchSize = 10
	for {
		// Получаем данные из шины
		msgs, err := r.nats.Sub.Fetch(batchSize)
		if err != nil {
			logrus.Errorf("[%s] error getting msgs, err: %v", handleName, err)
			time.Sleep(time.Duration(r.cfg.TimeSleepOnError) * time.Second)
			continue
		}

		var data [][]byte
		for _, msg := range msgs {
			data = append(data, msg.Data)
		}

		if err = handleFunc.Run(data); err != nil {
			logrus.Errorf("[%s] error handling msgs, err: %v", handleName, err)
			time.Sleep(time.Duration(r.cfg.TimeSleepOnError) * time.Second)
		}

		for i := range msgs {
			if err = msgs[i].Ack(); err != nil {
				logrus.Errorf("[%s] can not ack msgs: %v", handleName, err)
			}
		}

		time.Sleep(time.Duration(r.cfg.TimeSleepOnOk) * time.Second)
	}
}
