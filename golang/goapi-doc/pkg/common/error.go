package common

import (
	"net/http"
	"strings"
)

var (
	// ErrBodyParser json format error
	ErrBodyParser = NewBadRequestError("invalid json format")
	// ErrQueryParser query string format error
	ErrQueryParser = NewBadRequestError("invalid query string type")
	// ErrIdFormat id format error
	ErrIdFormat = NewBadRequestError("invalid id format")
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = NewNotFoundError("record not found")
	// ErrFromDatabase any error from database
	ErrFromDatabase = NewUnexpectedError("error from database")
	// ErrDbQuery error when select data from database
	ErrDbQuery = NewUnexpectedError("database query error")
	// ErrDbInsert error when insert data to database
	ErrDbInsert = NewUnexpectedError("database insert error")
	// ErrDbUpdate error when update data tp database
	ErrDbUpdate = NewUnexpectedError("database update error")
	// ErrDbDelete error when delete data from database
	ErrDbDelete = NewUnexpectedError("database delete error")
)

const (
	ErrMessageInvalidData = "invalid data see details"
)

type AppError struct {
	Code    int           `json:"code" example:"400"`
	Message string        `json:"message" example:"Invalid json body"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Target  string `json:"target" example:"Name"`
	Message string `json:"message" example:"Name field is required"`
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Is(target error) bool {
	return target.Error() == e.Message
}

func NewBadRequestError(message string) error {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewInvalidError(details string) error {
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: ErrMessageInvalidData,
		Details: parseError(details),
	}
}

func NewUnauthorizedError(message string) error {
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) error {
	return AppError{
		Code:    http.StatusForbidden,
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

func parseError(err string) []ErrorDetail {
	eds := []ErrorDetail{}
	errs := strings.Split(err, ",")
	for _, e := range errs {
		kv := strings.Split(e, ":")
		eds = append(eds, ErrorDetail{
			Target:  strings.TrimSpace(kv[0]),
			Message: strings.TrimSpace(kv[1]),
		})
	}
	return eds
}
