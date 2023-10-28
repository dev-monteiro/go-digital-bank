package configuration

import (
	"net/http"
)

type AppError struct {
	Message    string
	StatusCode int
}

func (appErr *AppError) Error() string {
	return appErr.Message
}

func NewUnknownError(err error) *AppError {
	return &AppError{
		Message:    "Unknown error: " + err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}
