package utils

import (
	"fmt"
	"net/http"
)

type Errors struct {
	Code     ErrorWrap
	Message  string
	Original error
}
type ErrorWrap int

const (
	ErrorPermission   ErrorWrap = http.StatusForbidden
	NotFound          ErrorWrap = http.StatusNotFound
	Internal          ErrorWrap = http.StatusInternalServerError
	ErrorBadRequest   ErrorWrap = http.StatusBadRequest
	ErrorUnauthorized ErrorWrap = http.StatusUnauthorized
)

func (e Errors) Error() string {
	return e.Message
}

func (e ErrorWrap) Wrap(eo error, msg string) Errors {
	return Errors{
		Code:     e,
		Message:  msg,
		Original: eo,
	}
}

func (e ErrorWrap) Wrapf(eo error, msg string, args ...interface{}) Errors {
	return Errors{
		Code:     e,
		Message:  fmt.Sprintf(msg, args...),
		Original: eo,
	}
}

func (e ErrorWrap) New(msg string) Errors {
	return Errors{
		Code:     e,
		Message:  msg,
		Original: nil,
	}
}

func (e ErrorWrap) Newf(msg string, args ...interface{}) Errors {
	return Errors{
		Code:     e,
		Message:  fmt.Sprintf(msg, args...),
		Original: nil,
	}
}
