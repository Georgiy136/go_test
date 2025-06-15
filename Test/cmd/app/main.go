package main

import (
	"github.com/sirupsen/logrus"
	"myapp/config"
	"myapp/internal/app"
)

//	@title			Swagger API
//	@version		1.0
//	@description	Swagger API for Golang Project

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

//	@BasePath	/

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
