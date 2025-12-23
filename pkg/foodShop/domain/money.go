package domain

import "fmt"

// Money stored as satang to avoid float issues
type Money int64

func THB(baht int64) Money { return Money(baht * 100) }

func (m Money) Add(x Money) Money { return m + x }
func (m Money) Sub(x Money) Money { return m - x }
func (m Money) MulInt(n int) Money { return m * Money(n) }

// Percent uses integer arithmetic => floor rounding automatically
func (m Money) Percent(p int64) Money {
	return Money(int64(m) * p / 100)
}

func (m Money) String() string {
	baht := int64(m) / 100
	satang := int64(m) % 100
	if satang < 0 {
		satang = -satang
	}
	return fmt.Sprintf("%d.%02d THB", baht, satang)
}
