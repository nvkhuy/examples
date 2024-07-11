package errs

import (
	"fmt"
	"net/http"
)

// Error error struct
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Header  int         `json:"header"`
	Detail  interface{} `json:"detail,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) WithMessagef(format string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) WithDetail(defail interface{}) *Error {
	e.Detail = defail
	return e
}

func (e *Error) WithDetailMessagef(format string, args ...interface{}) *Error {
	e.Detail = fmt.Sprintf(format, args...)
	return e
}

func New(code int, message string, header ...int) *Error {
	var hc = http.StatusBadRequest
	if len(header) > 0 && header[0] < 1000 {
		hc = header[0]
	}
	return &Error{Code: code, Message: message, Header: hc}
}
