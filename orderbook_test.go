package poloniex

import (
	"log"

	"github.com/k0kubun/pp"
)

func ExampleOrderBook() {
	p := NewPublicOnly()
	ob, err := p.OrderBook("BTC_ETH")
	if err != nil {
		log.Fatalln(err)
	}
	pp.Println(ob.Asks[0], ob.Bids[0])
}
