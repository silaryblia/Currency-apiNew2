package domain

import (
	"fmt"
	"strings"
)

func ValidateCreateCurrency(
	code CurrencyCode,
	rate Rate) error {
	if rate.Float64() <= 0 {
		return fmt.Errorf("rate must be greater than zero")
	}

	if err := ValidateCurrencyCode(code); err != nil {
		return err
	}
	return nil
}

func ValidateRateChange(oldRate, newRate Rate) error {
	if newRate.Float64() <= 0 {
		return fmt.Errorf("%w: rate must be greater than ZERO", ErrValidation)
	}
	return nil
}

func ValidateCurrencyCode(code CurrencyCode) error {
	s := string(code)
	if len(s) != 3 {
		return fmt.Errorf("currency code must be 3 characters")
	}

	s = strings.ToUpper(s)
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return fmt.Errorf("currency code must contain only letters")
		}
	}
	return nil
}
