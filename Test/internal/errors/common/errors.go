package common

import (
	"net/http"
)

type CommonError struct {
	Code     int
	Message  string
	HttpCode int
}

var (
	BadRequest                 = CommonError{Code: 1, Message: "errors.common.badRequest", HttpCode: http.StatusBadRequest}
	Unauthorized               = CommonError{Code: 2, Message: "errors.common.unauthorized", HttpCode: http.StatusUnauthorized}
	Forbidden                  = CommonError{Code: 3, Message: "errors.common.forbidden", HttpCode: http.StatusForbidden}
	NotFoundError              = CommonError{Code: 4, Message: "errors.common.notFound", HttpCode: http.StatusNotFound}
	InternalServerError        = CommonError{Code: 5, Message: "errors.common.internalServerError", HttpCode: http.StatusInternalServerError}
	ServiceUnavailable         = CommonError{Code: 6, Message: "errors.common.serviceUnavailable", HttpCode: http.StatusServiceUnavailable}
	ServiceUnprocessableEntity = CommonError{Code: 7, Message: "errors.common.unprocessableEntity", HttpCode: http.StatusUnprocessableEntity}
)

type BusinessError struct {
	Err     *CommonError
	Message string
}

func (err *BusinessError) Error() string {
	return err.Message
}

type DBError struct {
	Err     *CommonError
	Message string
}

func (err *DBError) Error() string {
	return err.Message
}
