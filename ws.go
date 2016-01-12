package poloniex

import "gopkg.in/beatgammit/turnpike.v2"

func (p *Poloniex) SubscribeTicker(handler turnpike.EventHandler) {
	p.InitWS()
	p.subscribedTo["ticker"] = true
	p.ws.Subscribe("ticker", handler)
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
