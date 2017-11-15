package poloniex

import (
	"fmt"
	"log"
)

func ExampleBalance() {
	p := New("config.json")
	balances, err := p.Balances()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", balances)
}
