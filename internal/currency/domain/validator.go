package domain

import "fmt"

func ValidateCreateCurrency(
	code CurrencyCode,
	rate Rate) error {
	// CurrencyCode и Rate уже валидны
	// оставляем расширяемость
	return nil
}

func ValidateRateChange(oldRate, newRate Rate) error {
	if newRate <= 0 {
		return fmt.Errorf("%w: rate must be greater than ZERO", ErrValidation)
	}
	return nil
}
