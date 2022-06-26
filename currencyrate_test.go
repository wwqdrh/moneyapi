package moneyapi

import (
	"fmt"
	"testing"
)

func TestXEMoneyAPI(t *testing.T) {
	api, err := NewxEAPI("Basic bG9kZXN0YXI6d2FHTDVTUXE2alQ1T0hRelVlS0pwSXNpNm1YcXQyWno=")
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
		fmt.Println(api.Rate("USD", item))
	}
}
