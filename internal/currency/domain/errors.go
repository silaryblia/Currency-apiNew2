package domain

import "errors"

var (
	ErrNotFound      = errors.New("currency not found")
	ErrAlreadyExists = errors.New("currency already exists")
	ErrInvalidArg    = errors.New("invalid argument")
)
