package domain

import (
	"errors"
	"strings"
)

func ValidateCurrency(code string, rate float64) error {
	code = strings.TrimSpace(code)

	if len(code) != 3 {
		return errors.New("currency code must be ISO-4217 (3 letters)")
	}

	if rate <= 0 {
		return errors.New("rate must be greater than zero")
	}

	return nil
}
