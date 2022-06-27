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

func init() {
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

func Currency2Flag(name string) string {
	val := currencyFlagMap[name]
	if val == "" {
		return "N/A"
	}
	return val
}
