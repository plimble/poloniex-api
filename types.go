package poloniex

import (
	"encoding/json"
	"strconv"
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

func F(i interface{}) float64 {
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0.0
		}
		return a
	case float64:
		return i
	case int64:
		return float64(i)
	case json.Number:
		a, err := i.Float64()
		if err != nil {
			return 0.0
		}
		return a
	}
	return 0.0
}
