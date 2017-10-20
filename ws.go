package poloniex

import (
	"encoding/json"

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

	WSOrder struct {
		Rate   float64 `json:"rate,string"`
		Type   string  `json:"type"`
		Amount float64 `json:"amount,string"`
	}
	WSTrade struct {
		TradeID string
		Rate    float64 `json:",string"`
		Amount  float64 `json:",string"`
		Type    string
		Date    string
	}
	WSOrderOrTrade []struct {
		Data json.RawMessage
		Type string
	}
)

func (p *Poloniex) SubscribeTicker(ch chan WSTicker) {
	p.InitWS()
	p.subscribedTo["ticker"] = true
	p.ws.Subscribe("ticker", p.makeTickerHandler(ch))
}

func (p *Poloniex) SubscribeOrder(code string, handler turnpike.EventHandler) {
	p.InitWS()
	p.subscribedTo[code] = true
	p.ws.Subscribe(code, handler)
}

func (p *Poloniex) SubscribeTrollbox(handler turnpike.EventHandler) {
	p.InitWS()
	p.subscribedTo["trollbox"] = true
	p.ws.Subscribe("trollbox", handler)
}

func (p *Poloniex) UnsubscribeTicker() {
	p.InitWS()
	p.Unsubscribe("ticker")
}

func (p *Poloniex) UnsubscribeOrder(code string) {
	p.InitWS()
	p.Unsubscribe(code)
}

func (p *Poloniex) UnsubscribeTrollbox() {
	p.InitWS()
	p.Unsubscribe("trollbox")
}

func (p *Poloniex) Unsubscribe(code string) {
	p.InitWS()
	if p.isSubscribed(code) {
		delete(p.subscribedTo, code)
		p.ws.Unsubscribe(code)
	}
}

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

func (p *Poloniex) makeOrderHandler(ch chan WSOrder) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
	}
}
