package app

import (
	"context"
	"fmt"
	"github.com/Georgiy136/go_test/web_service/config"
	"github.com/Georgiy136/go_test/web_service/internal/http"
	"github.com/Georgiy136/go_test/web_service/internal/http/middleware"
	nats_service "github.com/Georgiy136/go_test/web_service/internal/nats"
	cache "github.com/Georgiy136/go_test/web_service/internal/repo/redis"
	"github.com/Georgiy136/go_test/web_service/internal/usecase"
	db "github.com/Georgiy136/go_test/web_service/internal/usecase/repo/postgres"
	"github.com/Georgiy136/go_test/web_service/pkg/jaegerotel"
	nats_conn "github.com/Georgiy136/go_test/web_service/pkg/nats"
	"github.com/Georgiy136/go_test/web_service/pkg/postgres"
	"github.com/Georgiy136/go_test/web_service/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"time"
)

func Run(cfg *config.Config) {
	// Инициализация трейсера
	tp, err := jaegerotel.NewJaegerTracerProvider(
		viper.GetString("jaeger_host"),
		jaegerotel.WithConfig("diaplan_editor", viper.GetString("app_env")),
	)
	if err != nil {
		logrus.Fatalf("не удалось инициализировать jaeger tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = tp.Shutdown(ctx); err != nil {
			logrus.Fatalf("не удалось корректно остановить сервис jaeger tracer provider")
		}
	}(ctx)
	// END Инициализация трейсера

	// Старт трейсинга инийциализации сервисов
	tctx, span := jaegerotel.StartNewSpan("ServicesInitialization")

	// connections
	pg, err := postgres.NewPostgres(tctx, cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}

	redisConn, err := redis.NewConn(tctx, cfg.Redis)
	if err != nil {
		logrus.Errorf("app - Run - redis.New: %v", err)
	}

	nats, err := nats_conn.New(tctx, cfg.Nats)
	if err != nil {
		logrus.Errorf("app - Run - nats.New: %v", err)
	}

	// END Старт трейсинга инийциализации сервисов
	span.End()

	// очередь Nats для сохранения логов
	natsService := nats_service.NewNatsService(nats)

	// инициализация сервиса для сохранения логов в очередь
	logger := middleware.NewLogger(tctx, natsService)

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
