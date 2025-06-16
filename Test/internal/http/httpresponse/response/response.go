package response

import "myapp/internal/errors/common"

type (
	SuccessData struct {
		Data any `json:"data"`
	}

	ErrorData struct {
		Error any `json:"error"`
	}

	Error struct {
		Message     string      `json:"message"`
		Description string      `json:"description,omitempty"`
		Details     interface{} `json:"details"`
		Code        int         `json:"code"`
	}
)

func NewSuccess(data any) SuccessData {
	return SuccessData{
		Data: data,
	}
}

func NewError(err common.CommonError, description string, details interface{}) ErrorData {
	return ErrorData{
		Error: Error{
			Message:     err.Message,
			Description: description,
			Details:     details,
			Code:        err.Code,
		},
	}
}
