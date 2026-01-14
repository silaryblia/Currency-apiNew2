package domain

import "fmt"

type Rate float64

func NewRate(v float64) (Rate, error) {
	if v <= 0 {
		return 0, fmt.Errorf("rate must be greater than zero, got %f", v)
	}
	return Rate(v), nil
}

func (r Rate) Float64() float64 {
	return float64(r)
}
