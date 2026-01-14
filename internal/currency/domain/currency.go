package domain

import "time"

type Currency struct {
	Code     CurrencyCode
	Rate     Rate
	RateDate time.Time
}

func (c Currency) ToPrimitives() (string, float64, time.Time) {
	return c.Code.String(), c.Rate.Float64(), c.RateDate
}
