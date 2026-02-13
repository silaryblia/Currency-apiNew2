package http

import (
	"Currency-apiNew2/internal/currency/domain"
	"errors"
	"net/http"
)

func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
