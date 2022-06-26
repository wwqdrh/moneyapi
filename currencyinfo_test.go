package moneyapi

import (
	"fmt"
	"testing"
)

func TestCurrencyInfo(t *testing.T) {
	api, err := NewCurrencyInfo()
	if err != nil {
		t.Fatal(err)
	}

	target := []string{
		"USD", "JPY", "BGN", "CZK", "DKK", "GBP",
		"HUF", "PLN", "RON", "SEK", "CHF", "ISK",
		"NOK", "HRK", "TRY", "AUD", "BRL", "CAD",
		"CNY", "HKD", "IDR", "INR", "KRW", "MXN",
		"MYR", "NZD", "PHP", "SGD", "THB", "ZAR",
	}
	for _, item := range target {
		items, err := api.RelativeItem(item)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println(items)
		}
	}
}
