# Go Poloniex API wrapper
This API should be a complete wrapper for the [Poloniex api](https://poloniex.com/support/api/), including the public, private and websocket APIs.

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
	"fmt"
	"log"

	poloniex "github.com/pharrisee/poloniex-api"
)

var (
	p = poloniex.New("../config.json")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ch := make(chan poloniex.WSTicker)
	p.SubscribeTicker(ch)
	for oot := range ch {
		fmt.Printf("%+v", oot)
	}
}
```
