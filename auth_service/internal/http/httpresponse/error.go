package httpresponse

import (
	"errors"
	"github.com/Georgiy136/go_test/auth_service/internal/http/httpresponse/http_errors"
	"github.com/Georgiy136/go_test/auth_service/internal/http/httpresponse/response"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, details interface{}) {
	var customError *http_errors.CustomError
	if errors.As(err, &customError) {
		SendError(c, *customError.Err, customError.Description, details)
		return
	}
	SendErrorInternalServerError(c, err.Error(), nil)
}

func SendError(c *gin.Context, commonErr http_errors.CommonHttpError, description string, details interface{}) {
	c.AbortWithStatusJSON(commonErr.HttpCode, response.NewError(commonErr, description, details))
}

func SendFailBadRequest(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.BadRequest, description, details)
}

func SendFailUnauthorized(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.Unauthorized, description, details)
}

func SendFailForbidden(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.Forbidden, description, details)
}

func SendFailNotFound(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.NotFoundError, description, details)
}

func SendFailUnprocessableEntity(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.ServiceUnprocessableEntity, description, details)
}

func SendErrorInternalServerError(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.InternalServerError, description, details)
}

func SendErrorServiceUnavailable(c *gin.Context, description string, details interface{}) {
	SendError(c, http_errors.ServiceUnavailable, description, details)
}
