package httpresponse

import (
	"errors"
	"github.com/gin-gonic/gin"
	"myapp/internal/errors/common"
	"myapp/internal/http/httpresponse/response"
)

func HandleError(c *gin.Context, err error, description string, details interface{}) {
	var businessError *common.BusinessError // ошибка бизнес-логики
	if errors.As(err, &businessError) {
		if businessError.Err == nil {
			businessError.Err = &common.ServiceUnprocessableEntity
		}
		SendError(c, *businessError.Err, description, details)
		return
	}
	var dbError *common.DBError // ошибка БД
	if errors.As(err, &dbError) {
		if dbError.Err == nil {
			dbError.Err = &common.ServiceUnprocessableEntity
		}
		SendError(c, *dbError.Err, description, details)
		return
	}
	SendErrorInternalServerError(c, err.Error(), nil)
}

func SendError(c *gin.Context, commonErr common.CommonError, description string, details interface{}) {
	c.AbortWithStatusJSON(commonErr.HttpCode, response.NewError(commonErr, description, details))
}

func SendFailBadRequest(c *gin.Context, description string, details interface{}) {
	SendError(c, common.BadRequest, description, details)
}

func SendFailUnauthorized(c *gin.Context, description string, details interface{}) {
	SendError(c, common.Unauthorized, description, details)
}

func SendFailForbidden(c *gin.Context, description string, details interface{}) {
	SendError(c, common.Forbidden, description, details)
}

func SendFailNotFound(c *gin.Context, description string, details interface{}) {
	SendError(c, common.NotFoundError, description, details)
}

func SendFailUnprocessableEntity(c *gin.Context, description string, details interface{}) {
	SendError(c, common.ServiceUnprocessableEntity, description, details)
}

func SendErrorInternalServerError(c *gin.Context, description string, details interface{}) {
	SendError(c, common.InternalServerError, description, details)
}

func SendErrorServiceUnavailable(c *gin.Context, description string, details interface{}) {
	SendError(c, common.ServiceUnavailable, description, details)
}
