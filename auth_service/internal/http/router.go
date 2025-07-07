package http

import (
	"github.com/gin-gonic/gin"

	"github.com/Georgiy136/go_test/auth_service/internal/usecase"
)

func NewRouter(router *gin.Engine, os usecase.GoodsUseCases) {
	goodsHandlers := &GoodsHandler{
		us: os,
	}

	// Routers
	api := router.Group("/api")
	{
		goods := api.Group("/goods")
		{
			goods.GET("/list", goodsHandlers.ListGoods)
		}
		good := api.Group("/good")
		{
			good.POST("/create", goodsHandlers.PostGoods)
			good.PATCH("/update", goodsHandlers.UpdateGood)
			good.PATCH("/reprioritize", goodsHandlers.ReprioritizeGood)
			good.DELETE("/remove", goodsHandlers.DeleteGood)
		}
	}
}
