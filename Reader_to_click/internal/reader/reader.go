package reader

import (
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	nats_pkg "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
	"time"
)

type Reader struct {
	cfg        config.Cron
	nats       *nats_pkg.Nats
	handleFunc handleFunc
}

type handleFunc interface {
	Run(data [][]byte) error
}

func NewReader(cfg config.Cron, nats *nats_pkg.Nats, handleFunc handleFunc) *Reader {
	return &Reader{
		cfg:        cfg,
		nats:       nats,
		handleFunc: handleFunc,
	}
}

func (c *Reader) Start() {
	for {
		// Получаем данные из шины
		msgs, err := c.nats.Sub.Fetch(10)
		if err != nil {
			logrus.Errorf("error getting msgs, err: %v", err)
			time.Sleep(time.Duration(c.cfg.TimeSleepOnError) * time.Second)
			continue
		}

		var data [][]byte
		for _, msg := range msgs {
			data = append(data, msg.Data)
		}

		if err = c.handleFunc.Run(data); err != nil {
			logrus.Errorf("error handling msgs, err: %v", err)
			time.Sleep(time.Duration(c.cfg.TimeSleepOnError) * time.Second)
		}

		for i := range msgs {
			if err = msgs[i].Ack(); err != nil {
				logrus.Errorf("can not ack msgs: %v", err)
			}
		}

		time.Sleep(time.Duration(c.cfg.TimeSleepOnOk) * time.Second)
	}
}
