package main

import (
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/app"
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
