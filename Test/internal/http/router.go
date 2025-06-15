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
		good := api.Group("/good")
		{
			good.POST("/", goodsHandlers.PostGoods)
			good.GET("/", goodsHandlers.ListGoods)
			good.GET("/:id", goodsHandlers.GetOneGood)
			good.PUT("/:id", goodsHandlers.UpdateGood)
			good.DELETE("/:id", goodsHandlers.DeleteGood)
		}
	}
}
