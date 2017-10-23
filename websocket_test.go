package poloniex

import (
	"fmt"
)

func ExampleWSTicker() {
	p := NewWithCredentials("Key goes here", "secret goes here")
	ch := p.SubscribeTicker()
	for oot := range ch {
		fmt.Printf("%+v", oot)
	}
}
