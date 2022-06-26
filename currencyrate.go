package moneyapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

////////////////////
// 币种实时汇率
// 数据来源:
// 1、欧洲银行
// 2、xeapi
////////////////////

// 一个小时请求一次
// base usd
type xeRateAPI struct {
	url  string
	auth string

	currencyBase string // 基础币种
	currencyMap  map[string]float64
	preupdate    time.Time
}

func NewxEAPI(auth string) (*xeRateAPI, error) {
	api := &xeRateAPI{
		url:          "https://www.xe.com/api/protected/midmarket-converter",
		auth:         auth,
		currencyBase: "USD",
		currencyMap:  map[string]float64{},
		preupdate:    time.Now().Add(-2 * time.Hour),
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	} else {
		return api, nil
	}

}

func (a *xeRateAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

	req, err := http.NewRequest("GET", a.url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", a.auth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := a.parseRate(body); err != nil {
		return err
	}

	a.preupdate = time.Now()
	return nil
}

// usd -> other
func (a *xeRateAPI) parseRate(body []byte) error {
	// "timestamp", "rates": {name: value...}
	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	rates, ok := data["rates"].(map[string]interface{})
	if !ok {
		return errors.New("xeapi 响应数据结构异常")
	}

	for name, val := range rates {
		if v, ok := val.(float64); ok {
			a.currencyMap[name] = v
		}
	}
	return nil
}

func (a *xeRateAPI) Rate(base string, target string) (float64, error) {
	if base != "USD" {
		return 0, errors.New("TODO, now base must be USD")
	}
	if err := a.checkUpdate(); err != nil {
		return 0, err
	}

	return a.currencyMap[target], nil
}
