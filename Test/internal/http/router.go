package http

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myapp/docs"
	"myapp/internal/usecase"
)

func NewRouter(router *gin.Engine, os usecase.GoodsUseCases) {
	goodsHandlers := &GoodsHandler{
		us: os,
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			good.PATCH("/reprioritize", goodsHandlers.UpdateGood)
			good.DELETE("/remove", goodsHandlers.DeleteGood)
		}
	}
}
