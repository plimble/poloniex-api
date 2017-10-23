package poloniex

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"gopkg.in/beatgammit/turnpike.v2"
)

type (
	//Poloniex describes the API
	Poloniex struct {
		Key          string
		Secret       string
		ws           *turnpike.Client
		subscribedTo map[string]bool
		debug        bool
		nonce        int64
		mutex        sync.Mutex
	}
)

const (
	// PUBLICURI is the address of the public API on Poloniex
	PUBLICURI = "https://poloniex.com/public"
	// PRIVATEURI is the address of the public API on Poloniex
	PRIVATEURI = "https://poloniex.com/tradingApi"
)

//InitWS is an attempt to work around the shitty poloniex WS api connection timeouts
func (p *Poloniex) InitWS() {
	if p.ws != nil {
		return
	}
	err := retry(100, 3*time.Second, func() error {
		t := &tls.Config{InsecureSkipVerify: true}
		u := "wss://api.poloniex.com"
		c, err := turnpike.NewWebsocketClient(turnpike.JSON, u, t)
		if err != nil {
			log.Println(err)
			return errors.Wrap(err, "open of websocket connection to "+u+" failed")
		}
		_, err = c.JoinRealm("realm1", nil)
		if err != nil {
			log.Println(err)
			return errors.Wrap(err, "joining realm1 failed")
		}
		p.ws = c
		return nil
	})
	if err != nil {
		log.Fatalln(errors.Wrap(err, "retries exhausted, fatal."))
	}
	p.subscribedTo = map[string]bool{}

}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		log.Println("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func (p *Poloniex) isSubscribed(code string) bool {
	_, ok := p.subscribedTo[code]
	return ok
}

//Debug turns on debugmode, which basically dumps all responses from the poloniex API REST server
func (p *Poloniex) Debug() {
	p.debug = true
}

func (p *Poloniex) getNonce() string {
	p.nonce++
	return fmt.Sprintf("%d", p.nonce)
}

// NewWithCredentials allows to pass in the key and secret directly
func NewWithCredentials(key, secret string) *Poloniex {
	p := &Poloniex{}
	p.Key = key
	p.Secret = secret
	p.nonce = time.Now().UnixNano()
	p.mutex = sync.Mutex{}
	return p
}

// NewWithConfig is the replacement function for New, pass in a configfile to use
func NewWithConfig(configfile string) *Poloniex {
	p := map[string]string{}
	// we have a configfile
	b, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "reading "+configfile+" failed."))
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "unmarshal of config failed."))
	}
	return NewWithCredentials(p["key"], p["secret"])

}

// NewPublicOnly allows the use of the public and websocket api only
func NewPublicOnly() *Poloniex {
	p := &Poloniex{}
	p.nonce = time.Now().UnixNano()
	p.mutex = sync.Mutex{}
	return p
}

// New is the legacy way to create a new client, here just to maintain api
func New(configfile string) *Poloniex {
	return NewWithConfig(configfile)
}

func trace(s string) (string, time.Time) {
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	elapsed := time.Since(startTime)
	log.Printf("trace end: %s, elapsed %f secs\n", s, elapsed.Seconds())
}

func toDecimal(i interface{}) decimal.Decimal {
	maxFloat := decimal.NewFromFloat(math.MaxFloat64)
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return maxFloat
		}
		return decimal.NewFromFloat(a)
	case float64:
		return decimal.NewFromFloat(i)
	case int64:
		return decimal.NewFromFloat(float64(i))
	case json.Number:
		a, err := i.Float64()
		if err != nil {
			return maxFloat
		}
		return decimal.NewFromFloat(a)
	}
	return maxFloat
}
