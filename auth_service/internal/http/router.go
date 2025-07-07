package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Georgiy136/go_test/auth_service/internal/usecase"
)

func NewRouter(router *gin.Engine, os usecase.AuthService) {
	authHandlers := &GoodsHandler{
		us: os,
	}

	// Routers
	api := router.Group("/api")
	{
		goods := api.Group("/goods")
		{
			goods.GET("/list", authHandlers.ListGoods)
		}
		good := api.Group("/good")
		{
			good.POST("/create", authHandlers.PostGoods)
			good.PATCH("/update", authHandlers.UpdateGood)
			good.PATCH("/reprioritize", authHandlers.ReprioritizeGood)
			good.DELETE("/remove", authHandlers.DeleteGood)
		}
	}
}
