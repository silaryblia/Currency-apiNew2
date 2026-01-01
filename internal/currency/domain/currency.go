package domain

import "time"

type Currency struct {
	Code     string    `json:"code"`
	Rate     float64   `json:"rate"`
	RateDate time.Time `json:"rate_date"`
}
