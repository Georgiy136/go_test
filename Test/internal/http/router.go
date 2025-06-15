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
	operator := router.Group("/api")
	{
		operator.POST("/", operatorHandlers.PostOperator)
		operator.GET("/", operatorHandlers.GetAllOperators)
		operator.GET("/:id", operatorHandlers.GetOneOperator)
		operator.PUT("/:id", operatorHandlers.UpdateOperator)
		operator.DELETE("/:id", operatorHandlers.DeleteOperator)
	}
}
