package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"myapp/config"
	"myapp/internal/http"
	"myapp/internal/http/middleware"
	db "myapp/internal/repository/postgres"
	cache "myapp/internal/repository/redis"
	nats_service "myapp/internal/sevice/nats"
	"myapp/internal/usecase"
	nats_conn "myapp/pkg/nats"
	"myapp/pkg/postgres"
	"myapp/pkg/redis"
)

func Run(cfg *config.Config) {
	// connections
	pg, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.CloseConn()

	redisConn, err := redis.NewConn(cfg.Redis)
	if err != nil {
		logrus.Infof("app - Run - redis.New: %v", err)
	}

	//click, err := clickhouse.New(cfg.Clickhouse)
	//if err != nil {
	//	logrus.Fatalf("app - Run - clickhouse.New: %v", err)
	//}
	//defer click.CloseConn()

	nats, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Fatalf("app - Run - nats.New: %v", err)
	}
	defer nats.CloseConn()

	// очередь Nats для сохранения логов
	natsLogs := nats_service.NewNatsService(nats)

	// инициализация middleware для сохранения логов в очередь
	logger := middleware.NewLogger(natsLogs, cfg.Nats.ChannelName)

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}
	//if err = click.MigrateUpClickhouse(); err != nil {
	//	logrus.Fatalf("app - Run - MigrateUpClickhouse: %v", err)
	//}

	// repo
	goodsRedis := cache.NewGoodsRedis(redisConn)
	goodsRepository := db.NewGoodsRepo(pg)

	// Use case
	goodsUseCases := usecase.NewGoodsUsecases(goodsRepository, goodsRedis)

	// HTTP Server
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(logger.LoggingMiddleware())

	http.NewRouter(router, *goodsUseCases)

	router.Run(fmt.Sprintf(":%d", cfg.Http.Port))
}
