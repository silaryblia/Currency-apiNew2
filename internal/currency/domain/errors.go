package domain

import "errors"

var (
	ErrNotFound        = errors.New("currency not found")
	ErrAlreadyExists   = errors.New("currency already exists")
	ErrInvalidArg      = errors.New("invalid argument")
	ErrValidation      = errors.New("validation error")
	ErrServiceNotReady = errors.New("service not ready")
	ErrInternal        = errors.New("internal error")
)
