# Go Poloniex API wrapper
This API should be a complete wrapper for the Poloniex api (https://poloniex.com), including the public, private and websocket APIs.

To use create a copy of config-example.json and fill in your API key and secret.

```json
{
    "key":"put your key here",
    "secret":"put your secret here"
}
```

You can also pass your key/secret pair in code rather than creating a config.json.

# Examples

## Public API

```go
package main

import (
	"log"

	"github.com/k0kubun/pp"
	poloniex "gitlab.com/wmlph/poloniex-api"
)

func main() {
	p := poloniex.NewPublicOnly()
	ob, err := p.OrderBook("BTC_ETH")
	if err != nil {
		log.Fatalln(err)
	}
	pp.Println(ob.Asks[0], ob.Bids[0])
}
```

## Private API

```go
package main

import (
	"fmt"
	"log"

	"gitlab.com/wmlph/poloniex-api"
)

func main() {
	p := poloniex.New("config.json")
	balances, err := p.Balances()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", balances)
}
```

## Websocket API

```go
package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/k0kubun/pp"
	"gitlab.com/wmlph/poloniex-api"
)

var (
	mutex = sync.Mutex{}
)

func main() {
	p := poloniex.New("config.json")
	p.SubscribeTicker(tickerHandler)
	t := time.Tick(10 * time.Second)
	for _ = range t {

	}
}

func tickerHandler(p []interface{}, n map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	t := poloniex.WSTicker{
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
	pp.Println(t)
}

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
	}
	return 0.0
}
```