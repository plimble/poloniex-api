package main

import (
	"fmt"
	"log"

	poloniex "github.com/pharrisee/poloniex-api"
)

func main() {
	p := poloniex.NewWithCredentials("Key goes here", "secret goes here")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ch := make(chan poloniex.WSTicker)
	p.SubscribeTicker(ch)
	for oot := range ch {
		fmt.Printf("%+v", oot)
	}
}
