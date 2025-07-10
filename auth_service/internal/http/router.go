package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Georgiy136/go_test/auth_service/internal/service"
)

func NewRouter(router *gin.Engine, us service.AuthService) {
	authHandlers := &AuthHandler{
		us: us,
	}

	// Routers
	api := router.Group("/api")
	{
		good := api.Group("/auth")
		{
			good.GET("/issue_tokens", authHandlers.GetTokens)
			good.PUT("/update_tokens", authHandlers.UpdateTokens)
		}
	}
}
