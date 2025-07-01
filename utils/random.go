package utils

import (
	"math/rand"

	"github.com/shopspring/decimal"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min)
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return int64(RandomInt(500, 10000)) / 100
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "PEN", "GBP", "JPY"}
	return currencies[rand.Intn(len(currencies))]
}

func RandomDecimal(min, max float64) decimal.Decimal {
	val := min + rand.Float64()*(max-min)
	return decimal.NewFromFloat(val).Round(2)
}
