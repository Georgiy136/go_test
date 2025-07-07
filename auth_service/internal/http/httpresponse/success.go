package httpresponse

import (
	"github.com/Georgiy136/go_test/web_service/internal/http/httpresponse/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendSuccess(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, response.NewSuccess(data))
}

func SendNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func SendSuccessOK(c *gin.Context, data any) {
	SendSuccess(c, http.StatusOK, data)
}
