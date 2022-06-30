package moneyapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wwqdrh/logger"
)

////////////////////
// 币种实时汇率
// 数据来源:
// 1、欧洲银行
// 2、中国银行
// 3、美国银行
// 2、xeapi
////////////////////

type ICurrencyRate interface {
	checkUpdate() error
	CurrencyMap() map[string]float64
	Status() bool // 当前接口状态
	Rate(base string, target string) (float64, error)
}

////////////////////
// 中国银行每天早上10点更新
// 只有26个有效币种数据
// 每天请求一次
////////////////////
type cnBankAPI struct {
	url          string
	status       bool
	currencyBase string
	currencyMap  map[string]float64
	preupdate    time.Time
}

func NewCnBankAPI() (*cnBankAPI, error) {
	api := &cnBankAPI{
		url:          "https://www.boc.cn/sourcedb/whpj",
		currencyBase: "CNY",
		currencyMap:  map[string]float64{},
		preupdate:    time.Now().AddDate(0, 0, -1),
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *cnBankAPI) Status() bool {
	return a.status
}

func (a *cnBankAPI) CurrencyMap() map[string]float64 {
	if err := a.checkUpdate(); err != nil {
		logger.DefaultLogger.Warn(err.Error())
	}

	return a.currencyMap
}

func (a *cnBankAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

	// 最多重试5次
next:
	for times := 0; times < 5; times += 1 {
		a.status = false
		var url string
		// 最多十页 毕竟没有那么多国家
		for i := 0; i < 10; i++ {
			if i == 0 {
				url = a.url + "/index.html"
			} else {
				url = a.url + fmt.Sprintf("/index_%d.html", i)
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				goto next
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				goto next
			}
			if err := a.parseRate(resp.Body); err != nil {
				logger.DefaultLogger.Warn(err.Error())
			}
			resp.Body.Close()
		}
		a.preupdate = time.Now()
		a.status = true
		return nil // 成功解析
	}

	return errors.New("重试失败5次")
}

func (a *cnBankAPI) parseRate(body io.ReadCloser) error {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return err
	}

	doc.Find("body > div > div.BOC_main > div.publish > div:nth-child(3) > table > tbody > tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
		// 100外币可以兑换的人民币
		name := s.Find("td:nth-child(1)").Text()
		value := s.Find("td:nth-child(3)").Text()
		if value == "" {
			return
		}

		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			return
		}
		a.currencyMap[Currency2Flag(name)] = 100 / valueFloat
	})

	return nil
}

func (a *cnBankAPI) Rate(base string, target string) (float64, error) {
	if err := a.checkUpdate(); err != nil {
		return 0, err
	}

	val1, ok := a.currencyMap[base]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	val2, ok := a.currencyMap[target]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	return val2 / val1, nil
}

////////////////////
// 欧洲银行
// base eur
////////////////////
type europeBankAPI struct {
	url         string
	status      bool
	preupdate   time.Time
	currencyMap map[string]float64
	rateData    *europeRateRes
}

type europeRateRes struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Gesmes  string   `xml:"gesmes,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Subject string   `xml:"subject"`
	Sender  struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"Sender"`
	Cube struct {
		Text string `xml:",chardata"`
		Cube []struct {
			Text string `xml:",chardata"`
			Time string `xml:"time,attr"`
			Cube []struct {
				Text     string `xml:",chardata"`
				Currency string `xml:"currency,attr"`
				Rate     string `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}

func NewEuropeBankAPI() (*europeBankAPI, error) {
	api := &europeBankAPI{
		url:         "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml",
		preupdate:   time.Now().Add(-1 * time.Hour),
		rateData:    new(europeRateRes),
		currencyMap: map[string]float64{},
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *europeBankAPI) Status() bool {
	return a.status
}

func (a *europeBankAPI) CurrencyMap() map[string]float64 {
	if err := a.checkUpdate(); err != nil {
		logger.DefaultLogger.Warn(err.Error())
	}

	return a.currencyMap
}

func (a *europeBankAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

next:
	for times := 0; times < 5; times += 1 {
		a.status = false
		req, err := http.NewRequest("GET", a.url, nil)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			goto next
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			goto next
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			goto next
		}
		if err := a.parseRate(body); err != nil {
			logger.DefaultLogger.Warn(err.Error())
			goto next
		}

		a.preupdate = time.Now()
		a.status = true
		return nil
	}
	return errors.New("重试失败5次")
}

func (a *europeBankAPI) parseRate(body []byte) error {
	if err := xml.Unmarshal(body, &a.rateData); err != nil {
		return err
	}
	if len(a.rateData.Cube.Cube) == 0 {
		return errors.New("欧洲银行相应数据为空")
	}

	latestCube := a.rateData.Cube.Cube[0]
	for _, item := range latestCube.Cube {
		value, err := strconv.ParseFloat(item.Rate, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			continue
		}
		a.currencyMap[item.Currency] = value
	}
	return nil
}

func (a *europeBankAPI) Rate(base string, target string) (float64, error) {
	if err := a.checkUpdate(); err != nil {
		return 0, err
	}

	val1, ok := a.currencyMap[base]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	val2, ok := a.currencyMap[target]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	return val2 / val1, nil
}

////////////////////
// 美国银行
////////////////////
type americaBankAPI struct {
	url       string
	status    bool
	preupdate time.Time

	rateData    *americaRes
	currencyMap map[string]float64
}

type americaRes struct {
	Rate []struct {
		CurrencyId              string      `json:"currencyId"`
		CurrencySiteId          interface{} `json:"currencySiteId"`
		CurrencyCode            string      `json:"currencyCode"`
		CountryName             string      `json:"countryName"`
		CurrencyName            string      `json:"currencyName"`
		CurrencyNamePlural      string      `json:"currencyNamePlural"`
		CurrencyBuyRate         string      `json:"currencyBuyRate"` // 一外币买多少USD
		CurrencySellRate        string      `json:"currencySellRate"`
		CurrencySmallestDenom   float64     `json:"currencySmallestDenom"`
		CheckBuyRate            interface{} `json:"checkBuyRate"`
		CheckSellRate           interface{} `json:"checkSellRate"`
		ConversionChartType     interface{} `json:"conversionChartType"`
		Region                  interface{} `json:"region"`
		Rank                    interface{} `json:"rank"`
		DeleteDate              interface{} `json:"deleteDate"`
		IsResponseNull          interface{} `json:"isResponseNull"`
		SystemDate              string      `json:"systemDate"`
		CurrencyBuyRateNumeric  float64     `json:"currencyBuyRateNumeric"`
		CurrencySellRateNumeric float64     `json:"currencySellRateNumeric"`
		CheckBuyRateNumeric     float64     `json:"checkBuyRateNumeric"`
		CheckSellRateNumeric    float64     `json:"checkSellRateNumeric"`
	} `json:"rate"`
}

func NewAmericaBankAPI() (*americaBankAPI, error) {
	api := &americaBankAPI{
		url:         "https://www.bankofamerica.com/salesservices/foreign-exchange/get-currdetails?_=",
		preupdate:   time.Now().Add(-1 * time.Hour),
		rateData:    new(americaRes),
		currencyMap: map[string]float64{},
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *americaBankAPI) Status() bool {
	return a.status
}

func (a *americaBankAPI) CurrencyMap() map[string]float64 {
	if err := a.checkUpdate(); err != nil {
		logger.DefaultLogger.Warn(err.Error())
	}

	return a.currencyMap
}

func (a *americaBankAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

next:
	for times := 0; times < 5; times += 1 {
		a.status = false
		req, err := http.NewRequest("GET", a.url+fmt.Sprint(time.Now().UnixMicro()), nil)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())

			goto next
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())

			goto next
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())

			goto next
		}
		if err := a.parseRate(body); err != nil {
			logger.DefaultLogger.Warn(err.Error())
			goto next
		}
		a.preupdate = time.Now()
		a.status = true
		return nil
	}
	return errors.New("重试失败5次")
}

func (a *americaBankAPI) parseRate(body []byte) error {
	if err := json.Unmarshal(body, &a.rateData); err != nil {
		return err
	}

	for _, item := range a.rateData.Rate {
		value, err := strconv.ParseFloat(item.CurrencyBuyRate, 64)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			continue
		}
		a.currencyMap[item.CurrencyCode] = 1 / value
	}
	return nil
}

func (a *americaBankAPI) Rate(base string, target string) (float64, error) {
	if err := a.checkUpdate(); err != nil {
		return 0, err
	}

	val1, ok := a.currencyMap[base]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	val2, ok := a.currencyMap[target]
	if !ok {
		return 0, errors.New("未找到" + base + "数据")
	}
	return val2 / val1, nil
}

// 一个小时请求一次
// base usd
type xeRateAPI struct {
	url    string
	auth   string
	status bool

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
		preupdate:    time.Now().Add(-1 * time.Hour),
	}
	if err := api.checkUpdate(); err != nil {
		return nil, err
	} else {
		return api, nil
	}

}

func (a *xeRateAPI) Status() bool {
	return a.status
}

func (a *xeRateAPI) checkUpdate() error {
	if !time.Now().After(a.preupdate.Add(time.Hour)) {
		return nil
	}

	a.status = false
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
	a.status = true
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
