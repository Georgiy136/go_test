package cron

import (
	"encoding/json"
	"github.com/Georgiy136/go_test/Reader_to_click/config"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/clickhouse"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/models"
	nats_pkg "github.com/Georgiy136/go_test/Reader_to_click/pkg/nats"
	"github.com/sirupsen/logrus"
	"time"
)

type Cron struct {
	cfg        config.Cron
	nats       *nats_pkg.Nats
	clickhouse *clickhouse.Clickhouse
}

func NewCron(cfg config.Cron, nats *nats_pkg.Nats, clickhouse *clickhouse.Clickhouse) *Cron {
	return &Cron{
		cfg:        cfg,
		nats:       nats,
		clickhouse: clickhouse,
	}
}

func (c *Cron) Start() {
	for {
		// Получаем данные из шины
		msgs, err := c.nats.Sub.Fetch(10)
		if err != nil {
			logrus.Errorf("error getting msgs, err: %v", err)
			time.Sleep(time.Duration(c.cfg.TimeSleepOnError) * time.Second)
			continue
		}

		var logs []models.Log
		for _, msg := range msgs {
			var log models.Log
			if err = json.Unmarshal(msg.Data, &log); err != nil {
				logrus.Errorf("error unmarshalling log, data: %v, err: %v", msg.Data, err)
				time.Sleep(time.Duration(c.cfg.TimeSleepOnError) * time.Second)
			}
			logs = append(logs, log)
		}

		if err != nil {
			continue
		}

		if len(logs) == 0 {
			time.Sleep(time.Duration(c.cfg.TimeSleepOnOk) * time.Second)
			continue
		}

		// отправляем в клик
		if err = c.clickhouse.SaveLogsToClick(logs); err != nil {
			logrus.Errorf("error saving logs to clickhouse: %v", err)
		}

		for i := range msgs {
			if err = msgs[i].Ack(); err != nil {
				logrus.Errorf("can not ack msgs: %v", err)
			}
		}

		time.Sleep(time.Duration(c.cfg.TimeSleepOnOk) * time.Second)
	}
}
