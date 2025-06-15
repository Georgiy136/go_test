package response

import "myapp/internal/errors/common"

type (
	ResponseData struct {
		Data any `json:"data"`
	}

	Error struct {
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
		Code    int                    `json:"code"`
	}
)

func NewSuccess(data any) ResponseData {
	return ResponseData{
		Data: data,
	}
}

func NewError(err common.CommonError, details map[string]interface{}) ResponseData {
	return ResponseData{
		Data: Error{
			Code:    err.Code,
			Message: err.Message,
			Details: details,
		},
	}
}
