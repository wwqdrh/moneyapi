package moneyapi

import (
	_ "embed"
	"encoding/json"

	"github.com/wwqdrh/logger"
)

//go:embed convert.json
var convert []byte

var (
	ConvertInfo        convertInfo
	flagNameMap        = map[string][]string{} // 标识转[中文，英文]
	nameFlagMap        = map[string]string{}   // 中文，英文 转标识
	currencyCNRelatMap = map[string][]string{}
	currencyENRelatMap = map[string][]string{}
	countryRelatMap    = map[string]string{}
	currencyFlagMap    = map[string]string{}
)

var (
	relations = []map[string]bool{blrRela, koreanRela, thbRela, idrRela, rubRela}

	blrRela = map[string]bool{
		"巴西里亚尔": true,
		"巴西雷亚尔": true,
	}
	koreanRela = map[string]bool{
		"韩国元": true,
		"韩元":  true,
	}
	thbRela = map[string]bool{
		"泰铢":  true,
		"泰国铢": true,
	}
	idrRela = map[string]bool{
		"印度尼西亚卢比": true,
		"印尼卢比":    true,
	}
	rubRela = map[string]bool{
		"卢布":    true,
		"俄罗斯卢布": true,
	}
)

type convertInfo struct {
	Currencies struct {
		Zh map[string]struct {
			Name         string   `json:"name"`
			RelatedTerms []string `json:"relatedTerms"`
			NamePlural   string   `json:"name_plural"`
		} `json:"zh-CN"`
		En map[string]struct {
			Name         string   `json:"name"`
			RelatedTerms []string `json:"relatedTerms"`
			NamePlural   string   `json:"name_plural"`
		} `json:"en"`
	} `json:"currencies"`
	Countries struct {
		Zh map[string]struct {
			Name string `json:"name"`
		} `json:"zh-CN"`
		En map[string]struct {
			Name string `json:"name"`
		} `json:"en"`
	} `json:"countries"`
}

// 如果xeapi无法获取就使用convert.json
func init() {
	// api, err := NewCurrencyInfo()
	// if err != nil {
	// 	logger.DefaultLogger.Warn(err.Error())
	// 	if err := json.Unmarshal(convert, &ConvertInfo); err != nil {
	// 		logger.DefaultLogger.Fatal(err.Error())
	// 	}
	// } else {
	// 	ConvertInfo.Countries = api.info.PageProps.CommonI18nResources.Countries
	// 	ConvertInfo.Currencies = api.info.PageProps.CommonI18nResources.Currencies
	// }

	if err := json.Unmarshal(convert, &ConvertInfo); err != nil {
		logger.DefaultLogger.Fatal(err.Error())
	}
	// 标识转中文英文
	// 英文、中文转标识
	for key, item := range ConvertInfo.Countries.Zh {
		flagNameMap[key] = append(flagNameMap[key], item.Name)
		nameFlagMap[item.Name] = key
	}
	for key, item := range ConvertInfo.Countries.En {
		flagNameMap[key] = append(flagNameMap[key], item.Name)
		nameFlagMap[item.Name] = key
	}

	// 币种的映射关系
	// 国家使用的币种
	for key, item := range ConvertInfo.Currencies.Zh {
		currencyCNRelatMap[key] = item.RelatedTerms
		currencyFlagMap[item.Name] = key
		for _, item := range item.RelatedTerms {
			countryRelatMap[item] = key
		}
	}
	for key, item := range ConvertInfo.Currencies.En {
		currencyENRelatMap[key] = item.RelatedTerms
		for _, item := range item.RelatedTerms {
			countryRelatMap[item] = key
		}
	}

}

// 标识转国家名
func CountryName(name string) []string {
	return flagNameMap[name]
}

// 国家名转标识
func CountryEnName(name string) string {
	return nameFlagMap[name]
}

// 获取使用币种的国家
// 0 中文
// 1 英文
func Currency2Country(name string, flag uint) []string {
	if flag != 0 && flag != 1 {
		return nil
	}

	var res []string
	if flag == 0 {
		res = currencyCNRelatMap[name]
	} else if flag == 1 {
		res = currencyENRelatMap[name]
	}
	new := res[:0]
	for _, item := range res {
		if _, ok := nameFlagMap[item]; ok {
			new = append(new, item)
		}
	}
	return new
}

func Country2Currency(name string) string {
	return countryRelatMap[name]
}

// 币种转为标识
func Currency2Flag(name string) string {
	for name := range RelationCurrency(name) {
		if val := currencyFlagMap[name]; val != "" {
			return val
		}
	}
	return name
}

// 处理中文可能存在同音字的问题
func RelationCurrency(name string) map[string]bool {
	for _, item := range relations {
		if _, ok := item[name]; ok {
			return item
		}
	}
	return map[string]bool{name: true}
}
