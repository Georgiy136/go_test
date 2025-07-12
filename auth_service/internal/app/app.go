package app

import (
	"fmt"
	"github.com/Georgiy136/go_test/auth_service/client"
	"github.com/Georgiy136/go_test/auth_service/config"
	"github.com/Georgiy136/go_test/auth_service/internal/http"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	"github.com/Georgiy136/go_test/auth_service/internal/service/crypter"
	db "github.com/Georgiy136/go_test/auth_service/internal/service/repo/postgres"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/jwt"
	"github.com/Georgiy136/go_test/auth_service/internal/service/token_generate/tokens"
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

	// crypter
	crypt := crypter.NewCrypter(cfg.Crypter.SignedKey)

	// jwt_token_gen
	jwtGen := jwt.NewJwtTokenGenerateGolangJwtV5()

	// tokens generate service
	tokenGenerator := token_generate.NewIssueTokensService(
		tokens.NewRefreshToken(jwtGen, cfg.RefreshToken),
		tokens.NewAccessToken(jwtGen, crypt, cfg.AccessToken),
	)

	userInfoClient := client.NewUserInfoClient(cfg.UserInfoClient)

	notificationClient := client.NewNotificationClient(cfg.NotificationClient)

	// Service
	authService := service.NewAuthService(tokenGenerator, crypt, userInfoClient, notificationClient, authRepo)

	// HTTP Server
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	http.NewRouter(router, *authService)

	if err = router.Run(fmt.Sprintf(":%d", cfg.Http.Port)); err != nil {
		logrus.Fatalf("app - Run - router.Run: %v", err)
	}
}
