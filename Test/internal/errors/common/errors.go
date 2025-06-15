package common

import ()

type CommonError struct {
	Code    int
	Message string
}

var (
	BadRequest          = CommonError{Code: 1, Message: "errors.common.BadRequest"}
	Unauthorized        = CommonError{Code: 2, Message: "errors.common.Unauthorized"}
	Forbidden           = CommonError{Code: 3, Message: "errors.common.Forbidden"}
	NotFoundError       = CommonError{Code: 4, Message: "errors.common.notFound"}
	InternalServerError = CommonError{Code: 5, Message: "errors.common.InternalServerError"}
	ServiceUnavailable  = CommonError{Code: 6, Message: "errors.common.ServiceUnavailable"}
)
