package httpresponse

import (
	"github.com/gin-gonic/gin"
	"myapp/internal/http/httpresponse/response"
	"net/http"
)

func SendSuccess(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, response.NewSuccess(data))
}

func SendNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func SendSuccessOK(c *gin.Context) {
	SendSuccess(c, http.StatusOK, struct{}{})
}

func SendSuccessOKWithData(c *gin.Context, data any) {
	SendSuccess(c, http.StatusOK, data)
}
