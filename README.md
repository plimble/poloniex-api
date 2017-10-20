This API should be a complete wrapper for the Poloniex api (https://poloniex.com), including the public, private and websocket APIs.

To use create a copy of config-example.json and fill in your API key and secret.

See examples...

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
