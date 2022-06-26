package moneyapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

////////////////////
// 世界货币信息
// 数据来源:
// 1、xeapi
////////////////////

type xeInfoAPI struct {
	url string

	info      *xeInfoRes
	preupdate time.Time
}

type xeInfoRes struct {
	PageProps struct {
		CommonI18nResources struct {
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
		} `json:"commonI18nResources"`
	} `json:"pageProps"`
}

func NewCurrencyInfo() (*xeInfoAPI, error) {
	api := &xeInfoAPI{
		url:       "https://www.xe.com/_next/data/afewa80gr-eiMmprdYaFN/zh-CN/currencyconverter/convert.json",
		info:      new(xeInfoRes),
		preupdate: time.Now().Add(-2 * time.Hour),
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *xeInfoAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

	req, err := http.NewRequest("GET", a.url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := a.parseInfo(body); err != nil {
		return err
	}
	a.preupdate = time.Now()
	return nil
}

func (a *xeInfoAPI) parseInfo(body []byte) error {
	return json.Unmarshal(body, a.info)
}

// 给定一个币种，获取相关联的国家
func (a *xeInfoAPI) RelativeItem(item string) ([]string, error) {
	if err := a.checkUpdate(); err != nil {
		return nil, err
	}

	info := a.info.PageProps.CommonI18nResources.Currencies.Zh
	if val, ok := info[item]; !ok {
		return nil, errors.New(item + "不存在")
	} else {
		return val.RelatedTerms, nil
	}
}
