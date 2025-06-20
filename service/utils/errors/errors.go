package errors

import (
	"errors"
	"net/http"
)

type Error struct {
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	return e.Message
}

func New(message string, status int) *Error {
	return &Error{Message: message, StatusCode: status}
}

var (
	ErrBadRequest   = New("bad request", http.StatusBadRequest)
	ErrUnauthorized = New("unauthorized", http.StatusUnauthorized)
	ErrForbidden    = New("forbidden", http.StatusForbidden)
	ErrNotFound     = New("not found", http.StatusNotFound)
	ErrConflict     = New("conflict", http.StatusConflict)
	ErrInternal     = New("internal server error", http.StatusInternalServerError)
)

func WriteHTTPError(w http.ResponseWriter, err error) {
	var customError *Error

	if ok := errors.As(err, &customError); ok {
		http.Error(w, customError.Message, customError.StatusCode)
	} else {
		http.Error(w, ErrInternal.Message, ErrInternal.StatusCode)
	}
}
