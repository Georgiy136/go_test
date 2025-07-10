package app

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/crypter"
	"github.com/Georgiy136/go_test/auth_service/internal/http"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	db "github.com/Georgiy136/go_test/auth_service/internal/service/repo/postgres"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token"
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
	authRepo := db.NewAuthRepo(pg)

	// token generate service
	tokenGenerator := token.NewIssueTokensService(cfg.Tokens)

	// crypter
	crypt := crypter.NewCrypter(cfg.Crypter.SignedKey)

	// Service
	authService := service.NewAuthService(tokenGenerator, crypt, authRepo)

	// HTTP Server
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	http.NewRouter(router, *authService)

	if err = router.Run(fmt.Sprintf(":%d", cfg.Http.Port)); err != nil {
		logrus.Fatalf("app - Run - router.Run: %v", err)
	}
}
