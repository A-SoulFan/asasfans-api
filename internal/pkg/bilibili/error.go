package bilibili

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Error struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Request *resty.Request `json:"request"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewError(code int, message string, request *resty.Request) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Request: request,
	}
}
