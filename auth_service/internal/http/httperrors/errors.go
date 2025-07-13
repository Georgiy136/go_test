package httperrors

import (
	"net/http"
)

type CommonError struct {
	Code     int
	Message  string
	HttpCode int
}

var (
	BadRequest                 = CommonError{Code: 1, Message: "errors.templates.badRequest", HttpCode: http.StatusBadRequest}
	Unauthorized               = CommonError{Code: 2, Message: "errors.templates.unauthorized", HttpCode: http.StatusUnauthorized}
	Forbidden                  = CommonError{Code: 3, Message: "errors.templates.forbidden", HttpCode: http.StatusForbidden}
	NotFoundError              = CommonError{Code: 4, Message: "errors.templates.notFound", HttpCode: http.StatusNotFound}
	InternalServerError        = CommonError{Code: 5, Message: "errors.templates.internalServerError", HttpCode: http.StatusInternalServerError}
	ServiceUnavailable         = CommonError{Code: 6, Message: "errors.templates.serviceUnavailable", HttpCode: http.StatusServiceUnavailable}
	ServiceUnprocessableEntity = CommonError{Code: 7, Message: "errors.templates.unprocessableEntity", HttpCode: http.StatusUnprocessableEntity}
)

type CustomError struct {
	Err         *CommonError
	Description string
}

func (err *CustomError) Error() string {
	return err.Description
}
