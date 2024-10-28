package common

import "net/http"

type ErrorCode string

const (
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrUnauthorized ErrorCode = "UNAUTHORIZIED"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrServer       ErrorCode = "SERVER_ERROR"
)

func (e ErrorCode) StatusCode() int {
	switch e {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewError(code ErrorCode, message string, err error) error {
	return AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
