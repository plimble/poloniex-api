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
