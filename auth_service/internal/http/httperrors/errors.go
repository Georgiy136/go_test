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
	BadRequest                 = CommonError{Code: 1, Message: "httperrors.common.badRequest", HttpCode: http.StatusBadRequest}
	Unauthorized               = CommonError{Code: 2, Message: "httperrors.common.unauthorized", HttpCode: http.StatusUnauthorized}
	Forbidden                  = CommonError{Code: 3, Message: "httperrors.common.forbidden", HttpCode: http.StatusForbidden}
	NotFoundError              = CommonError{Code: 4, Message: "httperrors.common.notFound", HttpCode: http.StatusNotFound}
	InternalServerError        = CommonError{Code: 5, Message: "httperrors.common.internalServerError", HttpCode: http.StatusInternalServerError}
	ServiceUnavailable         = CommonError{Code: 6, Message: "httperrors.common.serviceUnavailable", HttpCode: http.StatusServiceUnavailable}
	ServiceUnprocessableEntity = CommonError{Code: 7, Message: "httperrors.common.unprocessableEntity", HttpCode: http.StatusUnprocessableEntity}
)

type CustomError struct {
	Err         *CommonError
	Description string
}

func (err *CustomError) Error() string {
	return err.Description
}
