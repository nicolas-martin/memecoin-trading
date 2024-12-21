package errors

import "fmt"

type ErrorCode string

const (
	ErrNotFound          ErrorCode = "NOT_FOUND"
	ErrInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrForbidden         ErrorCode = "FORBIDDEN"
	ErrInternalServer    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDuplicateEntry    ErrorCode = "DUPLICATE_ENTRY"
	ErrInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrNotFound
	}
	return false
}

func IsInvalidInput(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrInvalidInput
	}
	return false
}
