package poloniex

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

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
		Pair    string
		Event   string
		TradeID int64
		Type    string
		Rate    float64
		Amount  float64
		Total   float64
		TS      time.Time
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
				for _, v := range orderbook {
					p.Emit(v.Event, v).Emit(v.Pair, v).Emit(v.Pair+"-"+v.Event, v)
				}
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

func (p *Poloniex) parseOrderbook(raw []interface{}) ([]WSOrderbook, error) {
	trades := []WSOrderbook{}
	marketID := int64(toFloat(raw[0]))
	pair, ok := p.byID[marketID]
	if !ok {
		return trades, errors.New("cannot parse to orderbook - invalid marketID")
	}
	for _, _v := range raw[2].([]interface{}) {
		v := _v.([]interface{})
		trade := WSOrderbook{}
		trade.Pair = pair
		switch v[0].(string) {
		case "i":
		case "o":
			trade.Event = "modify"
			if t := toFloat(v[3]); t == 0.0 {
				trade.Event = "remove"
			}
			trade.Type = "ask"
			if t := toFloat(v[1]); t == 1.0 {
				trade.Type = "bid"
			}
			trade.Rate = toFloat(v[2])
			trade.Amount = toFloat(v[3])
			trade.TS = time.Now()
		case "t":
			trade.Event = "trade"
			trade.TradeID = int64(toFloat(raw[1]))
			trade.Type = "sell"
			if t := toFloat(v[2]); t == 1.0 {
				trade.Type = "buy"
			}
			trade.Rate = toFloat(v[3])
			trade.Amount = toFloat(v[4])
			trade.Total = trade.Rate * trade.Amount
			t := time.Unix(int64(toFloat(v[5])), 0)
			trade.TS = t
		default:
		}
		trades = append(trades, trade)
	}
	return trades, nil
}
