package httpresponse

import (
	"github.com/gin-gonic/gin"
	"myapp/internal/errors/common"
	"myapp/internal/http/httpresponse/response"
	"net/http"
)

func SendError(c *gin.Context, httpCode int, commonErr common.CommonError, details map[string]interface{}) {
	c.AbortWithStatusJSON(httpCode, response.NewError(commonErr, details))
}

func SendFailBadRequestWithData(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusBadRequest, common.NotFoundError, details)
}

func SendFailUnauthorized(c *gin.Context) {
	SendError(c, http.StatusUnauthorized, common.CommonError{}, nil)
}

func SendFailForbidden(c *gin.Context) {
	SendError(c, http.StatusForbidden, common.CommonError{}, nil)
}

func SendFailNotFound(c *gin.Context) {
	SendError(c, http.StatusNotFound, common.CommonError{}, nil)
}

func SendErrorInternalServerError(c *gin.Context) {
	SendError(c, http.StatusInternalServerError, common.CommonError{}, nil)
}

func SendErrorServiceUnavailable(c *gin.Context, message string) {
	SendError(c, http.StatusServiceUnavailable, common.CommonError{}, nil)
}
