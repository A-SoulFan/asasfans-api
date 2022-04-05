package apperrors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	ValidationError = 6
	AuthError       = 7
	ServiceError    = 8
	UnknownError    = 9
)

type AppError struct {
	err          error
	Code         int
	Message      string
	ResponseType int
}

func (e AppError) Error() string {
	str := fmt.Sprintf("Error code: %d message: %s", e.Code, e.Message)
	if e.err != nil {
		return errors.Wrap(e.err, str).Error()
	}

	return str
}

func (e *AppError) Wrap(err error) *AppError {
	e.err = err
	return e
}

func NewError(code int, message string) *AppError {
	return &AppError{
		err:          nil,
		Code:         code,
		Message:      message,
		ResponseType: UnknownError,
	}
}
