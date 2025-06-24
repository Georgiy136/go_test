package main

import (
	"github.com/sirupsen/logrus"
	"myapp/config"
	"myapp/internal/app"
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
