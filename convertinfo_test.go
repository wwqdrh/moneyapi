package moneyapi

import (
	"fmt"
	"testing"
)

func TestCountryName(t *testing.T) {
	target := []string{
		"USD", "JPY", "BGN", "CZK", "DKK", "GBP",
		"HUF", "PLN", "RON", "SEK", "CHF", "ISK",
		"NOK", "HRK", "TRY", "AUD", "BRL", "CAD",
		"CNY", "HKD", "IDR", "INR", "KRW", "MXN",
		"MYR", "NZD", "PHP", "SGD", "THB", "ZAR",
	}
	for _, item := range target {
		t := Currency2Country(item, 0)
		fmt.Println(t)
		fmt.Println(Currency2Country(item, 1))
		fmt.Println(Country2Currency(t[0]))
	}
}
