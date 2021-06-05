package httperror

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"` // httpStatus
}

func New(code int, msg string, status ...int) error {
	err := &Error{
		Code:    code,
		Message: msg,
		Status:  http.StatusOK,
	}
	if len(status) > 0 {
		err.Status = status[0]
	}
	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d message = %s status = %d", e.Code, e.Message, e.Status)
}

func StatusBadRequest(code int, message string) error {
	return New(code, message, http.StatusBadRequest)
}

func StatusUnauthorized(code int, message string) error {
	return New(code, message, http.StatusUnauthorized)
}

func StatusForbidden(code int, message string) error {
	return New(code, message, http.StatusForbidden)
}

func StatusNotFound(code int, message string) error {
	return New(code, message, http.StatusNotFound)
}

func StatusMethodNotAllowed(code int, message string) error {
	return New(code, message, http.StatusMethodNotAllowed)
}

func StatusRequestTimeout(code int, message string) error {
	return New(code, message, http.StatusRequestTimeout)
}

func StatusConflict(code int, message string) error {
	return New(code, message, http.StatusConflict)
}

func StatusInternalServerError(code int, message string) error {
	return New(code, message, http.StatusInternalServerError)
}

func FromError(err error) (*Error, bool) {
	if e := new(Error); errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

var _ error = (*Error)(nil)
