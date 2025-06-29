package service

import (
	"encoding/json"
	"fmt"
	"github.com/Georgiy136/go_test/Reader_to_click/internal/models"
	"github.com/Georgiy136/go_test/Reader_to_click/pkg/clickhouse"
	"github.com/sirupsen/logrus"
)

type SendLogsToClick struct {
	click *clickhouse.Clickhouse
}

func NewSendLogsToClick(click *clickhouse.Clickhouse) *SendLogsToClick {
	return &SendLogsToClick{
		click: click,
	}
}

func (s *SendLogsToClick) Run(data [][]byte) error {
	var (
		logs []models.Log
		err  error
	)
	for _, msg := range data {
		var log models.Log
		if err = json.Unmarshal(msg, &log); err != nil {
			return fmt.Errorf("error unmarshalling log, data: %v, err: %v", string(msg), err)
		}
		logs = append(logs, log)
	}

	if len(logs) == 0 {
		return nil
	}

	// отправляем в клик
	if err = s.SaveLogsToClick(logs); err != nil {
		logrus.Errorf("error saving logs to clickhouse: %v", err)
	}
	return nil
}

func (s *SendLogsToClick) SaveLogsToClick(logs []models.Log) error {
	const insertDataFormat = "Insert into %s values ($1)"

	if _, err := s.click.Conn.Exec(fmt.Sprintf(insertDataFormat, s.click.Cfg.Dbname), logs); err != nil {
		return fmt.Errorf("can not SaveLogsToClick, err: %w", err)
	}
	return nil
}
