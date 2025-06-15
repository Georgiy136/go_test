package response

import "myapp/internal/errors/common"

type (
	Success struct {
		Data any `json:"data"`
	}

	Err struct {
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
		Code    int                    `json:"code"`
	}
)

func NewSuccess(data any) Success {
	return Success{
		Data: data,
	}
}

func NewError(err common.CommonError, details map[string]interface{}) Err {
	return Err{
		Code:    err.Code,
		Message: err.Message,
		Details: details,
	}
}
