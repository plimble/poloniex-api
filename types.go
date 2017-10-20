package poloniex

import (
	"encoding/json"
	"strconv"
)

func f(i interface{}) float64 {
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
