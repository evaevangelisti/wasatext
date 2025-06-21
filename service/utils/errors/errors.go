package errors

import (
	"encoding/json"
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
	ErrBadRequest   = New("Bad request", http.StatusBadRequest)
	ErrUnauthorized = New("Unauthorized access", http.StatusUnauthorized)
	ErrForbidden    = New("Forbidden access", http.StatusForbidden)
	ErrNotFound     = New("Resource not found", http.StatusNotFound)
	ErrConflict     = New("Conflict error", http.StatusConflict)
	ErrInternal     = New("Internal server error", http.StatusInternalServerError)
)

func WriteHTTPError(w http.ResponseWriter, err error) {
	var customError *Error

	w.Header().Set("Content-Type", "application/json")

	if ok := errors.As(err, &customError); ok {
		w.WriteHeader(customError.StatusCode)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": customError.Message}); err != nil {
			http.Error(w, customError.Message, customError.StatusCode)
		}
	} else {
		w.WriteHeader(ErrInternal.StatusCode)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": ErrInternal.Message}); err != nil {
			http.Error(w, ErrInternal.Message, ErrInternal.StatusCode)
		}
	}
}
