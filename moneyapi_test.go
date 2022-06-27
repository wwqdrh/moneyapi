package moneyapi

import (
	"fmt"
	"testing"
)

func TestMoneyAPI(t *testing.T) {
	api := NewMoneyAPI()
	target := []string{
		"CNY", "USD", "JPY", "BGN", "CZK", "DKK", "GBP",
		"HUF", "PLN", "RON", "SEK", "CHF", "ISK",
		"NOK", "HRK", "TRY", "AUD", "BRL", "CAD",
		"CNY", "HKD", "IDR", "INR", "KRW", "MXN",
		"MYR", "NZD", "PHP", "SGD", "THB", "ZAR",
	}
	for j := 0; j < len(target); j++ {
		for i := 0; i < j; i++ {
			fmt.Printf("%s -> %s: %.4f\n", target[i], target[j], api.Rate(target[i], target[j]))
		}
	}
}
