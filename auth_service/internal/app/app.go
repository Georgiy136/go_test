package app

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/http"
	"github.com/Georgiy136/go_test/auth_service/internal/usecase"
	db "github.com/Georgiy136/go_test/auth_service/internal/usecase/repo/postgres"
	"github.com/Georgiy136/go_test/auth_service/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	// connections
	pg, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New: %v", err)
	}

	// Накатываем миграции
	if err = pg.MigrateUpPostgres(); err != nil {
		logrus.Fatalf("app - Run - MigrateUpPostgres: %v", err)
	}

	// repo
	goodsRepository := db.NewAuthRepo(pg)

	// Use case
	goodsUseCases := usecase.NewAuthService(goodsRepository)

	// HTTP Server
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	http.NewRouter(router, *goodsUseCases)

	if err = router.Run(fmt.Sprintf(":%d", cfg.Http.Port)); err != nil {
		logrus.Fatalf("app - Run - router.Run: %v", err)
	}
}
