package poloniex

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	turnpike "gopkg.in/beatgammit/turnpike.v2"
)

type (
	subscription struct {
		Command string `json:"command"`
		Channel string `json:"channel"`
	}
)

func (p *Poloniex) dial() error {
	if p.ws != nil {
		err := p.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
		if err != nil {
			p.ws = nil
			return p.dial()
		}
		return nil
	}
	c, _, err := websocket.DefaultDialer.Dial(apiURL, nil)
	if err != nil {
		return errors.Wrap(err, "websocket connection failed")
	}
	c.SetPongHandler(func(appData string) error {
		log.Println("PongHandler", appData)
		return nil
	})
	go func() {
		for _ = range time.Tick(10 * time.Second) {
			if p.ws != nil {
				err := p.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
				if err != nil {
					p.ws = nil
					err := p.dial()
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}()
	p.ws = c
	return nil
}

func (p *Poloniex) sendWSMessage(msg interface{}) error {
	p.dial()
	msgs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshalling WSmessage failed")
	}
	log.Println(string(msgs))

	err = p.ws.WriteMessage(websocket.TextMessage, msgs)
	if err != nil {
		return errors.Wrap(err, "sending WSmessage failed")
	}
	return nil
}

//makeTickerHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) messageHandler(ch chan WSTicker) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
		t := WSTicker{
			Pair:          p[0].(string),
			Last:          toFloat(p[1]),
			Ask:           toFloat(p[2]),
			Bid:           toFloat(p[3]),
			PercentChange: toFloat(p[4]) * 100.0,
			BaseVolume:    toFloat(p[5]),
			QuoteVolume:   toFloat(p[6]),
			IsFrozen:      toFloat(p[7]) != 0.0,
			DailyHigh:     toFloat(p[8]),
			DailyLow:      toFloat(p[9]),
		}
		ch <- t
	}
}

//makeOrderHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) makeOrderHandler(coin string, ch WSOrderOrTradeChan) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
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

					o.Data.TS = d
				}
				ootTmp = append(ootTmp, o)
			}
			o := WSOrderOrTrade{Orders: ootTmp}
			ch <- o
		}
	}
}
