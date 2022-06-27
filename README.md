<p align='center'>
  <pre style="float:left;">
 _   .-')                      .-') _     ('-.                   ('-.        _ (`-.            
( '.( OO )_                   ( OO ) )  _(  OO)                 ( OO ).-.   ( (OO  )           
 ,--.   ,--.) .-'),-----. ,--./ ,--,'  (,------.   ,--.   ,--.  / . --. /  _.`     \   ,-.-')  
 |   `.'   | ( OO'  .-.  '|   \ |  |\   |  .---'    \  `.'  /   | \-.  \  (__...--''   |  |OO) 
 |         | /   |  | |  ||    \|  | )  |  |      .-')     /  .-'-'  |  |  |  /  | |   |  |  \ 
 |  |'.'|  | \_) |  |\|  ||  .     |/  (|  '--.  (OO  \   /    \| |_.'  |  |  |_.' |   |  |(_/ 
 |  |   |  |   \ |  | |  ||  |\    |    |  .--'   |   /  /\_    |  .-.  |  |  .___.'  ,|  |_.' 
 |  |   |  |    `'  '-'  '|  | \   |    |  `---.  `-./  /.__)   |  | |  |  |  |      (_|  |    
 `--'   `--'      `-----' `--'  `--'    `------'    `--'        `--' `--'  `--'        `--'    
  </pre>
</p>

<p align='center'>
方便地<sup><em>MoneyAPI</em></sup> 金融相关
<br> 
</p>

<br>

## 背景

封装现有API或者自己收集数据，为金融相关的的查询提供工具

## 特性

- 🗂 货币实时汇率
- 🗂 货币信息查询
- 🗂 货币使用者查询

## 使用手册

```go
import "github.com/wwqdrh/moneyapi"

func main() {
  // 币种转换
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

  // 谁使用这个币
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

  // ...
}
```