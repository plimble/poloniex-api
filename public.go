package poloniex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/franela/goreq"
	"github.com/k0kubun/pp"
	"github.com/shopspring/decimal"
)

type (
	Ticker      map[string]TickerEntry
	TickerEntry struct {
		Last        decimal.Decimal `json:",string"`
		Ask         decimal.Decimal `json:"lowestAsk,string"`
		Bid         decimal.Decimal `json:"highestBid,string"`
		Change      decimal.Decimal `json:"percentChange,string"`
		BaseVolume  decimal.Decimal `json:"baseVolume,string"`
		QuoteVolume decimal.Decimal `json:"quoteVolume,string"`
		IsFrozen    int64           `json:"isFrozen,string"`
	}

	DailyVolume          map[string]DailyVolumeEntry
	DailyVolumeEntry     map[string]decimal.Decimal
	DailyVolumeTemp      map[string]interface{}
	DailyVolumeEntryTemp map[string]interface{}

	OrderBook struct {
		Asks     []Order
		Bids     []Order
		IsFrozen bool
	}
	Order struct {
		Rate   decimal.Decimal
		Amount decimal.Decimal
	}

	OrderBookTemp struct {
		Asks     []OrderTemp
		Bids     []OrderTemp
		IsFrozen interface{}
	}
	OrderTemp        []interface{}
	OrderBookAll     map[string]OrderBook
	OrderBookAllTemp map[string]OrderBookTemp

	TradeHistory      []TradeHistoryEntry
	TradeHistoryEntry struct {
		ID     int64 `json:"globalTradeID"`
		Date   string
		Type   string
		Rate   decimal.Decimal `json:",string"`
		Amount decimal.Decimal `json:",string"`
		Total  decimal.Decimal `json:",string"`
	}

	ChartData      []ChartDataEntry
	ChartDataEntry struct {
		Date            int64
		High            decimal.Decimal
		Low             decimal.Decimal
		Open            decimal.Decimal
		Close           decimal.Decimal
		Volume          decimal.Decimal
		QuoteVolume     decimal.Decimal
		WeightedAverage decimal.Decimal
	}

	Currencies map[string]Currency
	Currency   struct {
		Name           string
		TxFee          decimal.Decimal `json:",string"`
		MinConf        decimal.Decimal
		DepositAddress string
		Disabled       int64
		Delisted       int64
		Frozen         int64
	}

	LoanOrders struct {
		Offers  []LoanOrder
		Demands []LoanOrder
	}
	LoanOrder struct {
		Rate     decimal.Decimal `json:",string"`
		Amount   decimal.Decimal `json:",string"`
		RangeMin decimal.Decimal
		RangeMax decimal.Decimal
	}
)

func (p *Poloniex) Ticker() (ticker Ticker, err error) {
	err = p.public("returnTicker", nil, &ticker)
	return
}

func (p *Poloniex) DailyVolume() (dailyVolume DailyVolume, err error) {
	dvt := DailyVolumeTemp{}
	err = p.public("return24hVolume", nil, &dvt)
	if err != nil {
		return
	}
	dailyVolume = DailyVolume{}
	for k := range dvt {
		v := dvt[k]
		dve := DailyVolumeEntry{}
		switch i := v.(type) {
		default:
			v := i.(map[string]interface{})
			for kk, vv := range v {
				dve[kk] = ToDecimal(vv)
			}
			dailyVolume[k] = dve
		case string:
			//ignore anything that isn't a map
		}
	}
	return
}

func (p *Poloniex) OrderBook(pair string) (orderBook OrderBook, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("depth", "40")
	obt := OrderBookTemp{}
	err = p.public("returnOrderBook", params, &obt)
	if err != nil {
		return
	}
	orderBook = tempToOrderBook(obt)
	return
}

func (p *Poloniex) OrderBookAll() (orderBook OrderBookAll, err error) {
	params := url.Values{}
	params.Add("depth", "5")
	params.Add("currencyPair", "all")
	obt := OrderBookAllTemp{}
	err = p.public("returnOrderBook", params, &obt)
	if err != nil {
		return
	}
	orderBook = OrderBookAll{}
	for k, v := range obt {
		orderBook[k] = tempToOrderBook(v)
	}
	return
}

func (p *Poloniex) TradeHistory(in ...interface{}) (tradeHistory TradeHistory, err error) {
	pp.Println(in)
	params := url.Values{}
	params.Add("currencyPair", in[0].(string))
	if len(in) > 1 {
		// we have a start date
		params.Add("start", fmt.Sprintf("%d", in[1].(int64)))
	}
	if len(in) > 2 {
		// we have an end date
		params.Add("end", fmt.Sprintf("%d", in[2].(int64)))
	}
	err = p.public("returnTradeHistory", params, &tradeHistory)
	return
}

func (p *Poloniex) ChartData(pair string) (chartData ChartData, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("start", fmt.Sprintf("%d", time.Now().Add(-24*time.Hour).Unix()))
	params.Add("end", "9999999999")
	params.Add("period", "300")
	err = p.public("returnChartData", params, &chartData)
	return
}

func (p *Poloniex) ChartDataPeriod(pair string, start, end time.Time) (chartData ChartData, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))
	params.Add("period", "300")
	err = p.public("returnChartData", params, &chartData)
	return
}

func (p *Poloniex) ChartDataCurrent(pair string) (chartData ChartData, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("start", fmt.Sprintf("%d", time.Now().Add(-5*time.Minute).Unix()))
	params.Add("end", "9999999999")
	params.Add("period", "300")
	err = p.public("returnChartData", params, &chartData)
	return
}

func (p *Poloniex) Currencies() (currencies Currencies, err error) {
	err = p.public("returnCurrencies", nil, &currencies)
	return
}

func (p *Poloniex) LoanOrders(currency string) (loanOrders LoanOrders, err error) {
	params := url.Values{}
	params.Add("currency", currency)
	err = p.public("returnLoanOrders", params, &loanOrders)
	return
}

func tempToOrderBook(obt OrderBookTemp) (ob OrderBook) {
	asks := obt.Asks
	bids := obt.Bids
	ob.IsFrozen = obt.IsFrozen.(string) != "0"
	ob.Asks = []Order{}
	ob.Bids = []Order{}
	for k := range asks {
		v := asks[k]
		price := ToDecimal(v[0])
		amount := ToDecimal(v[1])
		o := Order{Rate: price, Amount: amount}
		ob.Asks = append(ob.Asks, o)
	}
	for k := range bids {
		v := bids[k]
		price := ToDecimal(v[0])
		amount := ToDecimal(v[1])
		o := Order{Rate: price, Amount: amount}
		ob.Bids = append(ob.Bids, o)
	}
	return
}

func (p *Poloniex) public(command string, params url.Values, retval interface{}) (err error) {
	if p.debug {
		defer un(trace("public: " + command))
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if params == nil {
		params = url.Values{}
	}
	params.Add("command", command)
	req := goreq.Request{Uri: PUBLICURI, QueryString: params, Timeout: 130 * time.Second}
	res, err := req.Do()
	if err != nil {
		return
	}
	if p.debug {
		pp.Println(res.Request.URL.String())
	}

	defer res.Body.Close()

	s, err := res.Body.ToString()
	if err != nil {
		return
	}
	if p.debug {
		pp.Println(s)
	}
	err = json.Unmarshal([]byte(s), retval)
	return
}
