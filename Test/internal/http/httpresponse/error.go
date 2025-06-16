package httpresponse

import (
	"github.com/gin-gonic/gin"
	"myapp/internal/errors/common"
	"myapp/internal/http/httpresponse/response"
	"net/http"
)

func SendError(c *gin.Context, httpCode int, commonErr common.CommonError, description *string, details interface{}) {
	c.AbortWithStatusJSON(httpCode, response.NewError(commonErr, *description, details))
}

func SendFailBadRequest(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusBadRequest, common.BadRequest, description, details)
}

func SendFailUnauthorized(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusUnauthorized, common.Unauthorized, description, details)
}

func SendFailForbidden(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusForbidden, common.Forbidden, description, details)
}

func SendFailNotFound(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusNotFound, common.NotFoundError, description, details)
}

func SendErrorInternalServerError(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusInternalServerError, common.InternalServerError, description, details)
}

func SendErrorServiceUnavailable(c *gin.Context, description *string, details interface{}) {
	SendError(c, http.StatusServiceUnavailable, common.ServiceUnavailable, description, details)
}
