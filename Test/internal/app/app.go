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
	// connections
	pg, err := postgres.New(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.CloseConn()

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatalf("app - Run - redis.New: %v", err)
	}
	defer rdb.Close()

	//click, err := clickhouse.New(cfg.Clickhouse)
	//if err != nil {
	//	logrus.Fatalf("app - Run - clickhouse.New: %v", err)
	//}
	//defer click.CloseConn()

	//nats, err := nats.New(cfg.Nats)
	//if err != nil {
	//	logrus.Fatalf("app - Run - nats.New: %v", err)
	//}
	//defer nats.CloseConn()

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}
	//if err = click.MigrateUpClickhouse(); err != nil {
	//	logrus.Fatalf("app - Run - MigrateUpClickhouse: %v", err)
	//}

	// repo
	redisGoods := repository.NewGoodsRedis(rdb)

	goodsRepository := repository.NewGoods(pg)

	// Use case
	goodsUseCases := usecase.NewGoodsUsecases(goodsRepository, redisGoods*repository.GoodsRedis)

	// HTTP Server
	router := gin.Default()

	http.NewRouter(router, *goodsUseCases)

	router.Run(fmt.Sprintf(":%d", cfg.Http.Port))
}
