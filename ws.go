package poloniex

import (
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
	}

	// WSTickerChan is a onduit through which WSTicker items are sent
	WSTickerChan chan WSTicker

	//WSTrade describes a trade, a new order, or an order update
	WSTrade struct {
		TradeID string
		Rate    float64 `json:",string"`
		Amount  float64 `json:",string"`
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

func (p *Poloniex) StartWS() {
	p.dial()
	go func() {
		for {
			_, message, err := p.ws.ReadMessage()
			if err != nil {
				p.dial()
				log.Println("read:", err)
				continue
			}
			p.Emit("message", string(message))
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
