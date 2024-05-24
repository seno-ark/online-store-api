package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidRequest = errors.New("invalid data")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrPermission     = errors.New("permission denied")
	ErrUnprocessable  = errors.New("unprocessable entity")
	ErrInternalServer = errors.New("internal server error")
)

func NewErrNotFound(message string) error {
	return fmt.Errorf("%s:%w", message, ErrNotFound)
}

func NewErrInvalidRequest(message string) error {
	return fmt.Errorf("%s:%w", message, ErrInvalidRequest)
}

func NewErrUnauthorized(message string) error {
	return fmt.Errorf("%s:%w", message, ErrUnauthorized)
}

func NewErrPermission(message string) error {
	return fmt.Errorf("%s:%w", message, ErrPermission)
}

func NewErrUnprocessable(message string) error {
	return fmt.Errorf("%s:%w", message, ErrUnprocessable)
}

func NewErrInternalServer(message string) error {
	return fmt.Errorf("%s:%w", message, ErrInternalServer)
}

func ErrStatusCode(err error) (int, string) {
	errMessage := strings.Split(err.Error(), ":")[0]

	switch {

	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, errMessage

	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest, errMessage

	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, errMessage

	case errors.Is(err, ErrPermission):
		return http.StatusForbidden, errMessage

	case errors.Is(err, ErrUnprocessable):
		return http.StatusUnprocessableEntity, errMessage

	case errors.Is(err, ErrInternalServer):
		return http.StatusInternalServerError, errMessage

	default:
		return http.StatusInternalServerError, "oops"

	}
}
