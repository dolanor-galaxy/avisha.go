// Package currency provides minimalistic currency primitives that maintain
// precision to the mill (1/100 of a cent).
package currency

import "fmt"

// Currency with precision up to 1/100 of a cent.
type Currency int64

const (
	Mill   Currency = 1
	Cent            = 100 * Mill
	Dollar          = 100 * Cent
)

func (c Currency) Mills() int64 {
	return int64(c)
}

func (c Currency) Cents() int64 {
	return int64(c) / 10
}

func (c Currency) Dollars() float64 {
	var (
		dollars = c / Dollar
		cents   = (c % Dollar) / Cent
	)
	return float64(dollars) + float64(cents)
}

func (c Currency) String() string {
	return fmt.Sprintf("$%.02f", c.Dollars())
}
