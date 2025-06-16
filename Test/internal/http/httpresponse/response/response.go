package response

import "myapp/internal/errors/common"

type (
	ResponseData struct {
		Data any `json:"data"`
	}

	Error struct {
		Message     string      `json:"message"`
		Description string      `json:"description,omitempty"`
		Details     interface{} `json:"details"`
		Code        int         `json:"code"`
	}
)

func NewSuccess(data any) ResponseData {
	return ResponseData{
		Data: data,
	}
}

func NewError(err common.CommonError, description string, details interface{}) ResponseData {
	return ResponseData{
		Data: Error{
			Message:     err.Message,
			Description: description,
			Details:     details,
			Code:        err.Code,
		},
	}
}
