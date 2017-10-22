package poloniex

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/k0kubun/pp"
	"github.com/shopspring/decimal"
	"gopkg.in/beatgammit/turnpike.v2"
)

type (
	//WSTicker describes a ticker item
	WSTicker struct {
		Pair          string
		Last          decimal.Decimal
		Ask           decimal.Decimal
		Bid           decimal.Decimal
		PercentChange decimal.Decimal
		BaseVolume    decimal.Decimal
		QuoteVolume   decimal.Decimal
		IsFrozen      bool
		DailyHigh     decimal.Decimal
		DailyLow      decimal.Decimal
	}

	// WSTickerChan is a onduit through which WSTicker items are sent
	WSTickerChan chan WSTicker

	//WSTrade describes a trade, a new order, or an order update
	WSTrade struct {
		TradeID string
		Rate    decimal.Decimal `json:",string"`
		Amount  decimal.Decimal `json:",string"`
		Type    string
		Date    string
		TS      time.Time
	}

	//WSOrderOrTrade is a slice of WSTrades with an indicator of the type (trade, new order, update order)
	WSOrderOrTrade struct {
		Seq    int64
		Orders WSOrders
	}

	WSOrders []struct {
		Data WSTrade
		Type string
	}

	// WSOrderOrTradeChan is a onduit through which WSTicker items are sent
	WSOrderOrTradeChan chan WSOrderOrTrade
)

const (
	//SENTINEL is used to mark items without a sequence number
	SENTINEL = int64(-1)
)

//SubscribeTicker subscribes to the ticker feed and returns a channel over which it will send updates
func (p *Poloniex) SubscribeTicker() WSTickerChan {
	p.InitWS()
	p.subscribedTo["ticker"] = true
	ch := make(WSTickerChan)
	p.ws.Subscribe("ticker", p.makeTickerHandler(ch))
	return ch
}

//SubscribeOrder subscribes to the order and trade feed and returns a channel over which it will send updates
func (p *Poloniex) SubscribeOrder(code string) WSOrderOrTradeChan {
	p.InitWS()
	p.subscribedTo[code] = true
	ch := make(WSOrderOrTradeChan)
	p.ws.Subscribe(code, p.makeOrderHandler(code, ch))
	return ch
}

//UnsubscribeTicker ... I think you can guess
func (p *Poloniex) UnsubscribeTicker() {
	p.InitWS()
	p.Unsubscribe("ticker")
}

//UnsubscribeOrder ... I think you can guess
func (p *Poloniex) UnsubscribeOrder(code string) {
	p.InitWS()
	p.Unsubscribe(code)
}

//Unsubscribe from the relevant feed
func (p *Poloniex) Unsubscribe(code string) {
	p.InitWS()
	if p.isSubscribed(code) {
		delete(p.subscribedTo, code)
		p.ws.Unsubscribe(code)
	}
}

//makeTickerHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) makeTickerHandler(ch chan WSTicker) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
		t := WSTicker{
			Pair:          p[0].(string),
			Last:          toDecimal(p[1]),
			Ask:           toDecimal(p[2]),
			Bid:           toDecimal(p[3]),
			PercentChange: toDecimal(p[4]).Mul(decimal.NewFromFloat(100.0)),
			BaseVolume:    toDecimal(p[5]),
			QuoteVolume:   toDecimal(p[6]),
			IsFrozen:      !toDecimal(p[7]).Equal(decimal.NewFromFloat(0.0)),
			DailyHigh:     toDecimal(p[8]),
			DailyLow:      toDecimal(p[9]),
		}
		ch <- t
	}
}

//makeOrderHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) makeOrderHandler(coin string, ch WSOrderOrTradeChan) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
		seq := SENTINEL
		if s, ok := n["seq"]; ok {
			if i, err := strconv.Atoi(s.(string)); err == nil {
				seq = int64(i)
			} else {
				seq = SENTINEL
			}
		}
		b, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}
		oot := WSOrders{}
		err = json.Unmarshal(b, &oot)
		if err != nil {
			log.Println(err)
			return
		}
		ootTmp := WSOrders{}
		for _, o := range oot {
			if o.Type == "newTrade" {
				pp.Println("Date:", o.Data.Date)
				d, err := time.Parse("2006-01-02 15:04:05", o.Data.Date)
				if err != nil {
					log.Println(err)
				}
				o.Data.TS = d
			}
			ootTmp = append(ootTmp, o)
		}
		o := WSOrderOrTrade{Seq: seq, Orders: ootTmp}
		ch <- o
	}
}
