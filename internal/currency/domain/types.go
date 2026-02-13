package domain

import (
	"fmt"
	"strings"
)

type CurrencyCode string

func ParseCurrencyCode(raw string) (CurrencyCode, error) {
	code := strings.ToUpper(strings.TrimSpace(raw))

	if len(code) != 3 {
		return "", fmt.Errorf("invalid currency code: %q", raw)
	}

	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return "", fmt.Errorf("invalid currency code: %q", raw)
		}
	}
	return CurrencyCode(code), nil
}

func (c CurrencyCode) String() string {
	return string(c)
}
