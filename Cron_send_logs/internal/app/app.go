package app

import (
	"fmt"
	"github.com/Georgiy136/go_test/go_test/Cron_send_logs/pkg/clickhouse"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"myapp/config"
	"myapp/internal/http"
	"myapp/internal/http/middleware"
	nats_service "myapp/internal/nats"
	db "myapp/internal/repository/postgres"
	cache "myapp/internal/repository/redis"
	"myapp/internal/usecase"
	nats_conn "myapp/pkg/nats"
	"myapp/pkg/postgres"
	"myapp/pkg/redis"
)

func Run(cfg *config.Config) {
	// connections
	_, err = clickhouse.New(cfg.Clickhouse)
	if err != nil {
		logrus.Errorf("app - Run - clickhouse.New: %v", err)
	}

	nats, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - nats.New: %v", err)
	}

	// очередь Nats для сохранения логов
	natsLogs := nats_service.NewNatsService(nats)

	// инициализация сервиса для сохранения логов в очередь
	logger := middleware.NewLogger(natsLogs, cfg.Nats.ChannelName)

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}
	//if err = click.MigrateUpClickhouse(); err != nil {
	//	logrus.Errorf("app - Run - MigrateUpClickhouse: %v", err)
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
