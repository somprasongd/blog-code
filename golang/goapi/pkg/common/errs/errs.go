package errs

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnexpectedError(message string) error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
