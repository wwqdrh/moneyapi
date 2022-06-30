package moneyapi

import (
	"fmt"
	"testing"
)

func TestCnBankRate(t *testing.T) {
	api, err := NewCnBankAPI()
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range api.currencyMap {
		fmt.Printf("%s %.4f\n", key, value)
	}
}

func TestAmericaBankRate(t *testing.T) {
	api, err := NewAmericaBankAPI()
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range api.currencyMap {
		fmt.Printf("%s %.4f\n", key, value)
	}
}

func TestEuropeBankRate(t *testing.T) {
	api, err := NewEuropeBankAPI()
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range api.currencyMap {
		fmt.Printf("%s %.4f\n", key, value)
	}
}

func TestXEMoneyAPI(t *testing.T) {
	api, err := NewxEAPI("Basic bG9kZXN0YXI6d2FHTDVTUXE2alQ1T0hRelVlS0pwSXNpNm1YcXQyWno=")
	if err != nil {
		fmt.Println(err.Error())
		t.Skip("xe api 失效")
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
