package poloniex

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

var (
	mutex sync.Mutex
)

func TestWSTicker(t *testing.T) {
	p := New("config.json")
	p.SubscribeTicker(tickerHandler)
	time.Sleep(5 * time.Second)
	p.UnsubscribeTicker()
}

func TestWSTrades(t *testing.T) {
	p := New("config.json")
	p.SubscribeOrder("BTC_FCT", orderHandler)
	time.Sleep(5 * time.Second)
	p.UnsubscribeOrder("BTC_FCT")
}

func TestWSTrollbox(t *testing.T) {
	p := New("config.json")
	p.SubscribeTrollbox(trollboxHandler)
	time.Sleep(30 * time.Second)
	p.UnsubscribeTrollbox()
}

func TestTicker(t *testing.T) {
	p := New("config.json")
	c, err := p.Ticker()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestDailyVolume(t *testing.T) {
	p := New("config.json")
	c, err := p.DailyVolume()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestOrderBook(t *testing.T) {
	p := New("config.json")
	c, err := p.OrderBook("BTC_FCT")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	fmt.Println("ASKS")
	for k, v := range c.Asks {
		fmt.Println(k, v)
	}
	fmt.Println("\nBIDS")
	for k, v := range c.Bids {
		fmt.Println(k, v)
	}

	fmt.Println("IsFrozen?", c.IsFrozen)
}

func TestOrderBookAll(t *testing.T) {
	p := New("config.json")
	c, err := p.OrderBookAll()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	for k, v := range c {
		fmt.Println("ASKS")
		for kk, vv := range v.Asks {
			fmt.Println(kk, vv)
		}
		fmt.Println("\nBIDS")
		for kk, vv := range v.Bids {
			fmt.Println(kk, vv)
		}

		fmt.Println(k, ": IsFrozen?", v.IsFrozen)
	}
}

func TestTradeHistory(t *testing.T) {
	p := New("config.json")
	c, err := p.TradeHistory("BTC_FCT")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestChartData(t *testing.T) {
	p := New("config.json")
	c, err := p.ChartData("BTC_FCT")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestCurrencies(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.Currencies()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestLoanOrders(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.LoanOrders("BTC")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestBalances(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.Balances()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestAddresses(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.Addresses()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestGenerateNewAddress(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.GenerateNewAddress("BTS")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestDepositsWithdrawals(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.DepositsWithdrawals()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestOpenOrders(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.OpenOrders("BTC_FCT")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestOpenOrdersAll(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.OpenOrdersAll()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestPrivateTradeHistory(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.PrivateTradeHistory("BTC_FCT")
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestLoanOffer(t *testing.T) {
	p := New("config.json")
	p.Debug()
	//currency string, amount float64, duration int, renew bool, lendingRate float64
	c, err := p.LoanOffer("DASH", 0.00117188, 2, false, 0.0599)
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestPrivateTradeHistoryAll(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.PrivateTradeHistoryAll()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestOpenLoanOffers(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.OpenLoanOffers()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func TestToggleAutoRenew(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.ToggleAutoRenew(13181666)
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}
func TestActiveLoans(t *testing.T) {
	p := New("config.json")
	p.Debug()
	c, err := p.ActiveLoans()
	if err != nil {
		t.Fail()
		log.Println(err)
		return
	}
	pp.Println(c)
}

func tickerHandler(p []interface{}, n map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	t := WSTicker{
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

func orderHandler(p []interface{}, n map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	for k := range p {
		v := p[k].(map[string]interface{})
		if v["type"].(string) == "newTrade" {
			pp.Println("NEWTRADE", v["data"])
		} else if v["type"].(string) == "orderBookModify" {
			pp.Println("ORDERBOOKMODIFY", v["data"])
		} else if v["type"].(string) == "orderBookRemove" {
			pp.Println("ORDERBOOKREMOVE", v["data"])
		} else {
			pp.Println(v)
		}
	}
}

func trollboxHandler(p []interface{}, n map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	pp.Println(p)
}

func init() {
	mutex = sync.Mutex{}
}
