package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"myapp/config"
	"myapp/internal/http"
	"myapp/internal/repository"
	"myapp/internal/usecase"
	"myapp/pkg/postgres"
	"myapp/pkg/redis"
)

func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.CloseConn()

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatalf("app - Run - redis.New: %v", err)
	}
	defer rdb.Close()

	redis := repository.NewRedis(rdb)
	logrus.Infof("app - Run - redis - %v", redis)

	// repo
	goodsRepository := repository.NewGoods(pg)

	// Use case
	goodsUseCases := usecase.NewGoodsUsecases(goodsRepository)

	// HTTP Server
	router := gin.Default()

	http.NewRouter(router, *goodsUseCases)

	router.Run(fmt.Sprintf(":%d", cfg.Http.Port))
}
