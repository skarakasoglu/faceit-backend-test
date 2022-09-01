package user

import (
	"faceit-backend-test/internal/apierr"
	"net/http"
)

func repositoryError(err error) apierr.ApiError {
	return apierr.ApiError{
		StatusCode: http.StatusInternalServerError,
		Code:       "1000",
		Message:    err.Error(),
		Data:       nil,
	}
}
