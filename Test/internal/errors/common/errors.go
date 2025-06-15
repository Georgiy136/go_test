package common

type CommonError struct {
	Code    int
	Message string
}

var (
	NotFoundError = CommonError{Code: 1, Message: "errors.common.notFound"}
)
