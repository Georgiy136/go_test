package httperrors

import (
	"net/http"
)

type CommonHttpError struct {
	Code     int
	Message  string
	HttpCode int
}

var (
	BadRequest                 = CommonHttpError{Code: 1, Message: "errors.templates.badRequest", HttpCode: http.StatusBadRequest}
	Unauthorized               = CommonHttpError{Code: 2, Message: "errors.templates.unauthorized", HttpCode: http.StatusUnauthorized}
	Forbidden                  = CommonHttpError{Code: 3, Message: "errors.templates.forbidden", HttpCode: http.StatusForbidden}
	NotFoundError              = CommonHttpError{Code: 4, Message: "errors.templates.notFound", HttpCode: http.StatusNotFound}
	InternalServerError        = CommonHttpError{Code: 5, Message: "errors.templates.internalServerError", HttpCode: http.StatusInternalServerError}
	ServiceUnavailable         = CommonHttpError{Code: 6, Message: "errors.templates.serviceUnavailable", HttpCode: http.StatusServiceUnavailable}
	ServiceUnprocessableEntity = CommonHttpError{Code: 7, Message: "errors.templates.unprocessableEntity", HttpCode: http.StatusUnprocessableEntity}
)

type CustomError struct {
	Err         *CommonHttpError
	Description string
}

func (err *CustomError) Error() string {
	return err.Description
}
