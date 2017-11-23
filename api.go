package poloniex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/chuckpreslar/emission"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	//Poloniex describes the API
	Poloniex struct {
		Key           string
		Secret        string
		ws            *websocket.Conn
		debug         bool
		nonce         int64
		mutex         sync.Mutex
		emitter       *emission.Emitter
		subscriptions map[string]bool
	}
)

const (
	// PUBLICURI is the address of the public API on Poloniex
	PUBLICURI = "https://poloniex.com/public"
	// PRIVATEURI is the address of the public API on Poloniex
	PRIVATEURI = "https://poloniex.com/tradingApi"
)

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
	p.emitter = emission.NewEmitter()
	p.subscriptions = map[string]bool{}

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
	p.emitter = emission.NewEmitter()
	p.subscriptions = map[string]bool{}
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

func toFloat(i interface{}) float64 {
	maxFloat := float64(math.MaxFloat64)
	switch i := i.(type) {
	case string:
		a, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return maxFloat
		}
		return a
	case float64:
		return i
	case int64:
		return float64(i)
	case json.Number:
		a, err := i.Float64()
		if err != nil {
			return maxFloat
		}
		return a
	}
	return maxFloat
}

func toString(i interface{}) string {
	switch i := i.(type) {
	case string:
		return i
	case float64:
		return fmt.Sprintf("%.8f", i)
	case int64:
		return fmt.Sprintf("%d", i)
	case json.Number:
		return i.String()
	}
	return ""
}
