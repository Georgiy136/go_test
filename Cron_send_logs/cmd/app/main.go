package main

import (
	"github.com/Georgiy136/go_test/Cron_send_logs/config"
	"github.com/Georgiy136/go_test/Cron_send_logs/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
