package poloniex

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

const (
	apiURL = "wss://api2.poloniex.com/"
)

var (
	_ChannelIDs = map[string]string{
		"trollbox":  "1001",
		"ticker":    "1002",
		"footer":    "1003",
		"heartbeat": "1010",
	}
)

type (
	//WSTicker describes a ticker item
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
		PairID        int64
	}

	WSOrderbook struct {
	}
)

func (p *Poloniex) StartWS() {
	go func() {
		for {
			_, raw, err := p.ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			var message []interface{}
			err = json.Unmarshal(raw, &message)
			if err != nil {
				log.Println(err)
				continue
			}
			chid := int64(message[0].(float64))
			chids := toString(chid)
			if chid > 100.0 && chid < 1000.0 {
				// it's an orderbook
				orderbook, err := p.parseOrderbook(message)
				if err != nil {
					log.Println(err)
					continue
				}
				p.Emit("orderbook", orderbook)
			} else if chids == _ChannelIDs["ticker"] {
				// it's a ticker
				ticker, err := p.parseTicker(message)
				if err != nil {
					log.Println(err)
					continue
				}
				p.Emit("ticker", ticker)
			}
		}
	}()
}

func (p *Poloniex) Subscribe(chid string) error {
	ticker, err := p.Ticker()
	if err != nil {
		return errors.Wrap(err, "getting ticker for subscribe failed")
	}

	if c, ok := _ChannelIDs[chid]; ok {
		chid = c
	} else if t, ok := ticker[chid]; ok {
		chid = strconv.Itoa(int(t.ID))
	} else {
		return errors.New("unrecognised channelid in subscribe")
	}

	p.subscriptions[chid] = true
	message := subscription{Command: "subscribe", Channel: chid}
	return p.sendWSMessage(message)
}

func (p *Poloniex) Unsubscribe(chid string) error {
	ticker, err := p.Ticker()
	if err != nil {
		return errors.Wrap(err, "getting ticker for unsubscribe failed")
	}

	if c, ok := _ChannelIDs[chid]; ok {
		chid = c
	} else if t, ok := ticker[chid]; ok {
		chid = strconv.Itoa(int(t.ID))
	} else {
		return errors.New("unrecognised channelid in unsubscribe")
	}
	message := subscription{Command: "subscribe", Channel: chid}
	delete(p.subscriptions, chid)
	return p.sendWSMessage(message)
}

func (p *Poloniex) parseTicker(raw []interface{}) (WSTicker, error) {
	wt := WSTicker{}
	pp.Println(raw)
	var rawInner []interface{}
	if len(raw) <= 2 {
		return wt, errors.New("cannot parse to ticker")
	}
	rawInner = raw[2].([]interface{})
	marketID := int64(toFloat(rawInner[0]))
	pair, ok := p.byID[marketID]
	if !ok {
		return wt, errors.New("cannot parse to ticker - invalid marketID")
	}

	wt.Pair = pair
	wt.PairID = marketID
	wt.Last = toFloat(rawInner[1])
	wt.Ask = toFloat(rawInner[2])
	wt.Bid = toFloat(rawInner[3])
	wt.PercentChange = toFloat(rawInner[4])
	wt.BaseVolume = toFloat(rawInner[5])
	wt.QuoteVolume = toFloat(rawInner[6])
	wt.IsFrozen = toFloat(rawInner[7]) != 0.0
	wt.DailyHigh = toFloat(rawInner[8])
	wt.DailyLow = toFloat(rawInner[9])

	return wt, nil
}

func (p *Poloniex) parseOrderbook(raw []interface{}) (WSOrderbook, error) {
	return WSOrderbook{}, nil
}
