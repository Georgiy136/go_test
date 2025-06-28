package app

import (
	"fmt"
	"github.com/Georgiy136/go_test/go_test/config"
	"github.com/Georgiy136/go_test/go_test/internal/http"
	"github.com/Georgiy136/go_test/go_test/internal/http/middleware"
	nats_service "github.com/Georgiy136/go_test/go_test/internal/nats"
	db "github.com/Georgiy136/go_test/go_test/internal/repository/postgres"
	cache "github.com/Georgiy136/go_test/go_test/internal/repository/redis"
	"github.com/Georgiy136/go_test/go_test/internal/usecase"
	nats_conn "github.com/Georgiy136/go_test/go_test/pkg/nats"
	"github.com/Georgiy136/go_test/go_test/pkg/postgres"
	"github.com/Georgiy136/go_test/go_test/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	// connections
	pg, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}

	redisConn, err := redis.NewConn(cfg.Redis)
	if err != nil {
		logrus.Errorf("app - Run - redis.New: %v", err)
	}

	nats, err := nats_conn.New(cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - cron.New: %v", err)
	}

	// очередь Nats для сохранения логов
	natsService := nats_service.NewNatsService(nats)

	// инициализация сервиса для сохранения логов в очередь
	logger := middleware.NewLogger(natsService)

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}

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

	if err = router.Run(fmt.Sprintf(":%d", cfg.Http.Port)); err != nil {
		logrus.Fatalf("app - Run - router.Run: %v", err)
	}
}
