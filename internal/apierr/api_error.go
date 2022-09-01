package apierr

import "net/http"

type ApiError struct {
	StatusCode int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (err ApiError) Error() string {
	return err.Message
}

func InternalServerError() ApiError {
	return ApiError{
		StatusCode: http.StatusInternalServerError,
		Code:       "0000",
		Message:    "internal server error",
		Data:       nil,
	}
}

func NotImplemented() ApiError {
	return ApiError{
		StatusCode: http.StatusInternalServerError,
		Code:       "0001",
		Message:    "not implemented",
		Data:       nil,
	}
}

func BadRequest(message string) ApiError {
	return ApiError{
		StatusCode: http.StatusBadRequest,
		Code:       "0002",
		Message:    message,
		Data:       nil,
	}
}
