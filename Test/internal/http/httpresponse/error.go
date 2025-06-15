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

func SendFailBadRequest(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusBadRequest, common.BadRequest, details)
}

func SendFailUnauthorized(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusUnauthorized, common.Unauthorized, details)
}

func SendFailForbidden(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusForbidden, common.Forbidden, details)
}

func SendFailNotFound(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusNotFound, common.NotFoundError, details)
}

func SendErrorInternalServerError(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusInternalServerError, common.InternalServerError, details)
}

func SendErrorServiceUnavailable(c *gin.Context, details map[string]interface{}) {
	SendError(c, http.StatusServiceUnavailable, common.ServiceUnavailable, details)
}
