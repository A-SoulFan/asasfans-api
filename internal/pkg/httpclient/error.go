package httpclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Error struct {
	Code    int    `json:"code"`
	Body    []byte `json:"body"`
	Request *resty.Request
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, body: %s", e.Code, e.Body)
}

func NewError(code int, body []byte, request *resty.Request) *Error {
	return &Error{
		Code:    code,
		Body:    body,
		Request: request,
	}
}
