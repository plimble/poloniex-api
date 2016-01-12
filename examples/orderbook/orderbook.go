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
