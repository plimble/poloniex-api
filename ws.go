package poloniex

import (
	"encoding/json"
	"log"
	"time"

	"github.com/k0kubun/pp"
	"gopkg.in/beatgammit/turnpike.v2"
)

type (
	WSTicker struct {
		Pair          string
		Last          float64
		Ask           float64
		Bid           float64
		PercentChange float64
		BaseVolume    float64
		QuoteVolume   float64
		IsFrozen      bool
		DailyHigh     float64
		DailyLow      float64
	}

	WSTrade struct {
		TradeID string
		Rate    float64 `json:",string"`
		Amount  float64 `json:",string"`
		Type    string
		Date    string
		TS      time.Time
	}
	WSOrderOrTrade []struct {
		Data WSTrade
		Type string
	}
)

//SubscribeTicker subscribes to the ticker feed
func (p *Poloniex) SubscribeTicker(ch chan WSTicker) {
	p.InitWS()
	p.subscribedTo["ticker"] = true
	p.ws.Subscribe("ticker", p.makeTickerHandler(ch))
}

//SubsribeOrder subscribes to the order and trade feed
func (p *Poloniex) SubscribeOrder(code string, ch chan WSOrderOrTrade) {
	p.InitWS()
	p.subscribedTo[code] = true
	p.ws.Subscribe(code, p.makeOrderHandler(code, ch))
}

//UnsubscribeTicker.... I think you can guess
func (p *Poloniex) UnsubscribeTicker() {
	p.InitWS()
	p.Unsubscribe("ticker")
}

//UnsubscribeOrder.... I think you can guess
func (p *Poloniex) UnsubscribeOrder(code string) {
	p.InitWS()
	p.Unsubscribe(code)
}

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
		ch <- t
	}
}

//makeOrderHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) makeOrderHandler(coin string, ch chan WSOrderOrTrade) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
		b, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}
		oot := WSOrderOrTrade{}
		err = json.Unmarshal(b, &oot)
		if err != nil {
			log.Println(err)
			return
		}
		ootTmp := WSOrderOrTrade{}
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
		ch <- ootTmp
	}
}
