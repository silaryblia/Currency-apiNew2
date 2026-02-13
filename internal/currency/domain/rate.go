package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Rate struct {
	decimal.Decimal
}

func NewRate(value interface{}) (Rate, error) {
	var d decimal.Decimal
	var err error

	switch v := value.(type) {
	case float64:
		d = decimal.NewFromFloat(v).Round(4)
	case string:
		d, err = decimal.NewFromString(v)
		if err == nil {
			d = d.Round(4)
		}
	case int:
		d = decimal.NewFromInt(int64(v))
	case int64:
		d = decimal.NewFromInt(v)
	case decimal.Decimal:
		d = v.Round(4)
	default:
		return Rate{}, fmt.Errorf("unsupported type for Rate: %T", value)
	}

	if err != nil {
		return Rate{}, fmt.Errorf("invalid rate: %w", err)
	}

	if d.IsNegative() {
		return Rate{}, fmt.Errorf("rate cannot be negative")
	}

	return Rate{Decimal: d}, nil
}

func (r Rate) Float64() float64 {
	f, _ := r.Decimal.Float64()
	return f
}

func MustRate(value interface{}) Rate {
	rate, err := NewRate(value)
	if err != nil {
		panic(err)
	}
	return rate
}

func RateFromFloat64(f float64) Rate {
	return Rate{Decimal: decimal.NewFromFloat(f).Round(4)}
}

func RateFromString(s string) (Rate, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Rate{}, err
	}
	return Rate{Decimal: d.Round(4)}, err
}

func CBRRate(s string) (Rate, error) {
	s = replaceCommandWithDot(s)
	return RateFromString(s)
}

func replaceCommandWithDot(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r == ',' {
			runes[i] = '.'
		}
	}
	return string(runes)
}

func (r Rate) IsPositive() bool {
	return r.GreaterThan(decimal.Zero)
}

func (r Rate) Diff(other Rate) decimal.Decimal {
	if other.IsZero() {
		return decimal.Zero
	}

	return r.
		Sub(other.Decimal).
		Abs().
		Div(other.Decimal)
}
