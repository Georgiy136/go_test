package http

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myapp/docs"
	"myapp/internal/usecase"
)

func NewRouter(router *gin.Engine, os usecase.OperatorUseCases) {
	operatorHandlers := &OperatorHandler{
		us: os,
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	api := router.Group("/api")
	{
		good := api.Group("/good")
		{
			good.POST("/", operatorHandlers.PostOperator)
			good.GET("/", operatorHandlers.GetAllOperators)
			good.GET("/:id", operatorHandlers.GetOneOperator)
			good.PUT("/:id", operatorHandlers.UpdateOperator)
			good.DELETE("/:id", operatorHandlers.DeleteOperator)
		}
	}
}
