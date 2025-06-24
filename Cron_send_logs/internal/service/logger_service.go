package usecase

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/Cron_send_logs/config"
	"github.com/Georgiy136/go_test/Cron_send_logs/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

type LogService struct {
	cfg             config.Cron
	timeSleepPeriod time.Duration
}

func NewLogService(cfg config.Cron) *LogService {
	return &LogService{
		cfg: cfg,
	}
}

func (us *LogService) SendLogsToClickhouse() {
	createdGood, err := us.db.CreateGoods(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("GoodUseCases - AddGoods - us.db.CreateGoods: %w", err)
	}

	return createdGood, nil
}
