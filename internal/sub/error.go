package sub

import "faceit-backend-test/internal/apierr"

func SubscribeError(message string) apierr.ApiError {
	return apierr.ApiError{
		StatusCode: 500,
		Code:       "2000",
		Message:    message,
		Data:       nil,
	}
}
