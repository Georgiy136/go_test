package common

type CommonError struct {
	Code    int
	Message string
}

var (
	BadRequest                 = CommonError{Code: 1, Message: "errors.common.badRequest"}
	Unauthorized               = CommonError{Code: 2, Message: "errors.common.unauthorized"}
	Forbidden                  = CommonError{Code: 3, Message: "errors.common.forbidden"}
	NotFoundError              = CommonError{Code: 4, Message: "errors.common.notFound"}
	InternalServerError        = CommonError{Code: 5, Message: "errors.common.internalServerError"}
	ServiceUnavailable         = CommonError{Code: 6, Message: "errors.common.serviceUnavailable"}
	ServiceUnprocessableEntity = CommonError{Code: 7, Message: "errors.common.unprocessableEntity"}
)

//const ProcedureError = "RAISE EXCEPTION procedure error"

type BusinessError struct {
	Message string
}

func (err BusinessError) Error() string {
	return err.Message
}
