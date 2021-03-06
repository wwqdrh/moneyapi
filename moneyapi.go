package moneyapi

import (
	"sync"
	"time"

	"github.com/wwqdrh/logger"
)

// graph base on bank rate info
type MoneyAPI struct {
	cnRateAPI ICurrencyRate
	europAPI  ICurrencyRate
	amerAPI   ICurrencyRate

	// 有向带权图
	graph      map[string]map[string]float64
	lastUpdate time.Time
}

func NewMoneyAPI() *MoneyAPI {
	api := &MoneyAPI{
		graph: map[string]map[string]float64{
			"CNY": {},
			"USD": {},
			"EUR": {},
		},
		lastUpdate: time.Now().Add(-1 * time.Hour),
	}
	api.Update()
	return api
}

func (a *MoneyAPI) Update2() {
	if !time.Now().After(a.lastUpdate.Add(1 * time.Hour)) {
		return
	}

	wait := sync.WaitGroup{}
	wait.Add(3)
	go func() {
		defer wait.Done()
		if a.cnRateAPI == nil || !a.cnRateAPI.Status() {
			cnRate, err := NewCnBankAPI()
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
			}
			a.cnRateAPI = cnRate
		}
	}()
	go func() {
		defer wait.Done()
		if a.europAPI == nil || !a.europAPI.Status() {
			europRate, err := NewEuropeBankAPI()
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
			}
			a.europAPI = europRate
		}
	}()

	go func() {
		defer wait.Done()
		if a.amerAPI == nil || !a.amerAPI.Status() {
			amerRate, err := NewAmericaBankAPI()
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
			}
			a.amerAPI = amerRate
		}
	}()
	wait.Wait()

	for item, val := range a.cnRateAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["CNY"][item] = val
		a.graph[item]["CNY"] = 1 / val
	}

	for item, val := range a.europAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["EUR"][item] = val
		a.graph[item]["EUR"] = 1 / val
	}

	for item, val := range a.amerAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["USD"][item] = val
		a.graph[item]["USD"] = 1 / val
	}
	a.lastUpdate = time.Now()
}

func (a *MoneyAPI) Update() {
	if a.cnRateAPI == nil || !a.cnRateAPI.Status() {
		cnRate, err := NewCnBankAPI()
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
		a.cnRateAPI = cnRate
	}

	if a.europAPI == nil || !a.europAPI.Status() {
		europRate, err := NewEuropeBankAPI()
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
		a.europAPI = europRate
	}

	if a.amerAPI == nil || !a.amerAPI.Status() {
		amerRate, err := NewAmericaBankAPI()
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
		}
		a.amerAPI = amerRate
	}

	for item, val := range a.cnRateAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["CNY"][item] = val
		a.graph[item]["CNY"] = 1 / val
	}

	for item, val := range a.europAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["EUR"][item] = val
		a.graph[item]["EUR"] = 1 / val
	}

	for item, val := range a.amerAPI.CurrencyMap() {
		if a.graph[item] == nil {
			a.graph[item] = map[string]float64{}
		}
		a.graph["USD"][item] = val
		a.graph[item]["USD"] = 1 / val
	}
}

// 在graph中寻找从base到target的路径
// bfs 寻找最短路径
// 1 base -> ? target
func (a *MoneyAPI) Rate(base, target string) float64 {
	a.Update()

	type node struct {
		Name  string
		Value float64 // 从1base到这个节点所等价的价值 但是这个浮点型导致如果直接使用*node在比较的时候会失败
	}

	root := &node{base, 1}
	visit := map[string]*node{root.Name: root}
	queue := []*node{root}
	for len(queue) > 0 {
		length := len(queue)
		for i := 0; i < length; i++ {
			cur := queue[0]
			queue = queue[1:]

			for next, val := range a.graph[cur.Name] {
				// 1 cur = val next
				nextNode := &node{
					next, val * cur.Value,
				}
				if next == target {
					return nextNode.Value
				}

				if _, ok := visit[nextNode.Name]; ok {
					continue
				}

				queue = append(queue, nextNode)
				visit[nextNode.Name] = nextNode
			}
		}
	}
	return -1
}
