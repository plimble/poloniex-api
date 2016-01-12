package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/k0kubun/pp"
	"gitlab.com/wmlph/poloniex-api"
)

var (
	mutex = sync.Mutex{}
)

func main() {
	p := poloniex.New("config.json")
	p.SubscribeTicker(tickerHandler)
	t := time.Tick(10 * time.Second)
	for _ = range t {

	}
}

func tickerHandler(p []interface{}, n map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	t := poloniex.WSTicker{
		Pair:          p[0].(string),
		Last:          f(p[1]),
		Ask:           f(p[2]),
		Bid:           f(p[3]),
		PercentChange: f(p[4]) * 100.0,
		BaseVolume:    f(p[5]),
		QuoteVolume:   f(p[6]),
		IsFrozen:      p[7].(float64) != 0.0,
		DailyHigh:     f(p[8]),
		DailyLow:      f(p[9]),
	}
	pp.Println(t)
}

func f(i interface{}) float64 {
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0.0
		}
		return a
	case float64:
		return i
	}
	return 0.0
}
