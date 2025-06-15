package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"myapp/config"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/usecase"
	"myapp/pkg/postgres"
	"myapp/pkg/redis"
)

func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.CloseConn()

	// Накатываем миграции
	pg.MigrateUpPostgres()

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer rdb.Close()

	redis := repository.NewRedis(rdb)
	logrus.Debugf("app - Run - redis - %v", redis)

	operatorRepository := repository.NewOperator(pg)

	// Use case
	operatorUseCases := usecase.NewOperatorUsecases(operatorRepository)

	// HTTP Server
	router := gin.Default()

	handler.NewRouter(router, *operatorUseCases)

	router.Run(fmt.Sprintf(":%d", cfg.Http.Port))
}
