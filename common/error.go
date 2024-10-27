package common

type ErrorCode string

const (
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrorUnauthorized ErrorCode = "UNAUTHORIZIED"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrServer         ErrorCode = "SERVER_ERROR"
)

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
