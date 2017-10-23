package poloniex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/franela/goreq"
	"github.com/shopspring/decimal"
)

type (
	//Base is the common fields returned by a call tothepoloniex API
	Base struct {
		Error    string
		Success  int64
		Response string
	}

	//Balances are the complete balance map returned by the poloniex API
	Balances map[string]Balance
	//Balance is a single balance entry used in the Balances map
	Balance struct {
		Available decimal.Decimal `json:",string"`
		OnOrders  decimal.Decimal `json:"onOrders,string"`
		BTCValue  decimal.Decimal `json:"btcValue,string"`
	}

	accountBalancesTemp struct {
		Exchange map[string]string
		Margin   map[string]string
		Lending  map[string]string
	}

	//Account holds the balances in the various wallet accounts
	AccountBalances struct {
		Exchange map[string]decimal.Decimal
		Margin   map[string]decimal.Decimal
		Lending  map[string]decimal.Decimal
	}

	//Addresses holds the various deposit addresses foreach coin
	Addresses map[string]string

	//DepositsWithdrawals holds the history of deposit and withdrawal
	DepositsWithdrawals struct {
		Deposits    []deposit
		Withdrawals []withdrawal
	}
	deposit struct {
		Currency      string
		Address       string
		Amount        decimal.Decimal `json:",string"`
		Confirmations int64
		TXID          string `json:"txid"`
		Timestamp     int64
		Status        string
	}
	withdrawal struct {
		WithdrawalNumber int64 `json:"withdrawalNumber"`
		Currency         string
		Address          string
		Amount           decimal.Decimal `json:",string"`
		Timestamp        int64
		Status           string
	}

	//OpenOrders is the list of open orders for the pair specified
	OpenOrders []OpenOrder
	//OpenOrder is a singular entry used in the OpenOrders type
	OpenOrder struct {
		OrderNumber int64 `json:",string"`
		Type        string
		Rate        decimal.Decimal `json:",string"`
		Amount      decimal.Decimal `json:",string"`
		Total       decimal.Decimal `json:",string"`
	}
	//OpenOrdersAll is used for all pairs
	OpenOrdersAll map[string]OpenOrders

	PrivateTradeHistory      []PrivateTradeHistoryEntry
	PrivateTradeHistoryEntry struct {
		Date        string
		Rate        decimal.Decimal `json:",string"`
		Amount      decimal.Decimal `json:",string"`
		Total       decimal.Decimal `json:",string"`
		OrderNumber int64           `json:",string"`
		Type        string
	}
	PrivateTradeHistoryAll map[string]PrivateTradeHistory

	OrderTrades []OrderTrade
	OrderTrade  struct {
		GlobalTradeID int64           `json:"globalTradeID"`
		TradeID       int64           `json:"tradeID"`
		CurrencyPair  string          `json:"currencyPair"`
		Type          string          `json:"type"`
		Rate          decimal.Decimal `json:"rate,string"`
		Amount        decimal.Decimal `json:"amount,string"`
		Total         decimal.Decimal `json:"total,string"`
		Fee           decimal.Decimal `json:"fee,string"`
		Date          string          `json:"date"`
	}

	Buy struct {
		OrderNumber int64 `json:",string"`
		// ResultingTrades []ResultingTrade
	}
	ResultingTrade struct {
		Amount  decimal.Decimal `json:",string"`
		Rate    decimal.Decimal `json:",string"`
		Date    string
		Total   decimal.Decimal `json:",string"`
		TradeID string          `json:"tradeID"`
		Type    string
	}
	Sell struct {
		Buy
	}

	MoveOrder struct {
		Base
		OrderNumber int64 `json:",string"`
		// ResultingTrades []ResultingTrade
	}

	Withdraw struct {
		Base
	}

	FeeInfo struct {
		MakerFee        decimal.Decimal `json:"makerFee,string"`
		TakerFee        decimal.Decimal `json:"takerFee,string"`
		ThirtyDayVolume decimal.Decimal `json:"thirtyDayVolume,string"`
		NextTier        decimal.Decimal `json:"nextTier,string"`
	}

	AvailableAccountBalances struct {
		Exchange map[string]decimal.Decimal
		Margin   map[string]decimal.Decimal
		Lending  map[string]decimal.Decimal
	}
	AvailableAccountBalancesTemp struct {
		Exchange map[string]json.Number
		Margin   map[string]json.Number
		Lending  map[string]json.Number
	}

	TradableBalances map[string]TradableBalance
	TradableBalance  map[string]decimal.Decimal

	TradableBalancesTemp map[string]TradableBalanceTemp
	TradableBalanceTemp  map[string]json.Number

	TransferBalance struct {
		Base
		Message string `json:"message"`
	}

	MarginAccountSummary struct {
		TotalValue         decimal.Decimal `json:"totalValue,string"`
		ProfitLoss         decimal.Decimal `json:"pl,string"`
		LendingFees        decimal.Decimal `json:"lendingFees,string"`
		NetValue           decimal.Decimal `json:"netValue,string"`
		TotalBorrowedValue decimal.Decimal `json:"totalBorrowedValue,string"`
		CurrentMargin      decimal.Decimal `json:"currentMargin,string"`
	}

	LoanOffer struct {
		Base
		OrderID int64 `json:"orderID"`
	}

	OpenLoanOffers map[string][]OpenLoanOffer
	OpenLoanOffer  struct {
		ID        int64           `json:"id"`
		Rate      decimal.Decimal `json:",string"`
		Amount    decimal.Decimal `json:",string"`
		Duration  int64
		Renewable bool
		AutoRenew int64 `json:"autoRenew"`
		Date      string
		DateTaken time.Time
	}

	ActiveLoans struct {
		Provided []ActiveLoan
	}
	ActiveLoan struct {
		ID        int64 `json:"id"`
		Currency  string
		Rate      decimal.Decimal `json:",string"`
		Amount    decimal.Decimal `json:",string"`
		Range     int64
		Renewable bool
		AutoRenew int64 `json:"autoRenew"`
		Date      string
		DateTaken time.Time
		Fees      decimal.Decimal `json:",string"`
	}
)

func (p *Poloniex) Balances() (balances Balances, err error) {
	p.private("returnCompleteBalances", nil, &balances)
	return
}

func (p *Poloniex) AccountBalances() (balances AccountBalances, err error) {
	b := accountBalancesTemp{}
	p.private("returnAvailableAccountBalances", nil, &b)
	balances = AccountBalances{Exchange: map[string]decimal.Decimal{}, Margin: map[string]decimal.Decimal{}, Lending: map[string]decimal.Decimal{}}
	for k, v := range b.Exchange {
		balances.Exchange[k], _ = decimal.NewFromString(v)
	}
	for k, v := range b.Margin {
		balances.Margin[k], _ = decimal.NewFromString(v)
	}
	for k, v := range b.Lending {
		balances.Lending[k], _ = decimal.NewFromString(v)
	}
	return
}

func (p *Poloniex) Addresses() (addresses Addresses, err error) {
	p.private("returnDepositAddresses", nil, &addresses)
	return
}

func (p *Poloniex) GenerateNewAddress(currency string) (address string, err error) {
	params := url.Values{}
	params.Add("currency", currency)
	b := Base{}
	err = p.private("generateNewAddress", params, &b)
	address = b.Response
	return
}

func (p *Poloniex) DepositsWithdrawals() (depositsWithdrawals DepositsWithdrawals, err error) {
	params := url.Values{}
	params.Add("start", fmt.Sprintf("%d", time.Now().Add(-5208*time.Hour).Unix()))
	params.Add("end", "9999999999")
	err = p.private("returnDepositsWithdrawals", params, &depositsWithdrawals)
	return
}

func (p *Poloniex) OpenOrders(pair string) (openOrders OpenOrders, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	err = p.private("returnOpenOrders", params, &openOrders)
	return
}

func (p *Poloniex) OpenOrdersAll() (openOrders OpenOrdersAll, err error) {
	params := url.Values{}
	params.Add("currencyPair", "all")
	err = p.private("returnOpenOrders", params, &openOrders)
	return
}

func (p *Poloniex) PrivateTradeHistory(pair string) (history PrivateTradeHistory, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	err = p.private("returnTradeHistory", params, &history)
	return
}

func (p *Poloniex) PrivateTradeHistoryAll() (history PrivateTradeHistoryAll, err error) {
	params := url.Values{}
	params.Add("currencyPair", "all")
	err = p.private("returnTradeHistory", params, &history)
	return
}

func (p *Poloniex) OrderTrades(orderNumber int64) (ot OrderTrades, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	err = p.private("returnOrderTrades", params, &ot)
	return
}

func (p *Poloniex) CancelOrder(orderNumber int64) (success bool, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	b := Base{}
	err = p.private("cancelOrder", params, &b)
	success = b.Success == 1
	return
}

func (p *Poloniex) Buy(pair string, rate, amount decimal.Decimal) (buy Buy, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	params.Add("amount", fmt.Sprintf("%.8f", amount))
	err = p.private("buy", params, &buy)
	return
}

func (p *Poloniex) BuyPostOnly(pair string, rate, amount decimal.Decimal) (buy Buy, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	params.Add("amount", fmt.Sprintf("%.8f", amount))
	params.Add("postOnly", "1")
	err = p.private("buy", params, &buy)
	return
}

func (p *Poloniex) Sell(pair string, rate, amount decimal.Decimal) (sell Sell, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	params.Add("amount", fmt.Sprintf("%.8f", amount))
	err = p.private("sell", params, &sell)
	return
}

func (p *Poloniex) SellPostOnly(pair string, rate, amount decimal.Decimal) (sell Sell, err error) {
	params := url.Values{}
	params.Add("currencyPair", pair)
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	params.Add("amount", fmt.Sprintf("%.8f", amount))
	params.Add("postOnly", "1")
	err = p.private("sell", params, &sell)
	return
}

func (p *Poloniex) Move(orderNumber int64, rate decimal.Decimal) (moveOrder MoveOrder, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	err = p.private("moveOrder", params, &moveOrder)
	return
}

func (p *Poloniex) MovePostOnly(orderNumber int64, rate decimal.Decimal) (moveOrder MoveOrder, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	params.Add("rate", fmt.Sprintf("%.8f", rate))
	err = p.private("moveOrder", params, &moveOrder)
	return
}

func (p *Poloniex) Withdraw(currency string, amount decimal.Decimal, address string) (w Withdraw, err error) {
	params := url.Values{}
	params.Add("currency", currency)
	params.Add("amount", fmt.Sprintf("%f", amount))
	params.Add("address", address)
	err = p.private("withdraw", params, w)
	return
}

func (p *Poloniex) FeeInfo() (fi FeeInfo, err error) {
	err = p.private("returnFeeInfo", nil, &fi)
	return
}

func (p *Poloniex) AvailableAccountBalances() (aab AvailableAccountBalances, err error) {
	aabt := AvailableAccountBalancesTemp{}
	err = p.private("returnAvailableAccountBalances", nil, &aabt)
	if err != nil {
		return
	}
	aab.Exchange = map[string]decimal.Decimal{}
	aab.Margin = map[string]decimal.Decimal{}
	aab.Lending = map[string]decimal.Decimal{}
	for k, v := range aabt.Exchange {
		aab.Exchange[k] = toDecimal(v)
	}
	for k, v := range aabt.Margin {
		aab.Margin[k] = toDecimal(v)
	}
	for k, v := range aabt.Lending {
		aab.Lending[k] = toDecimal(v)
	}
	return
}

func (p *Poloniex) TradableBalances() (tb TradableBalances, err error) {
	tbt := TradableBalancesTemp{}
	err = p.private("returnTradableBalances", nil, &tbt)
	if err != nil {
		return
	}
	tb = TradableBalances{}
	for k, v := range tbt {
		tb[k] = TradableBalance{}
		for kk, vv := range v {
			tb[k][kk] = toDecimal(vv)
		}
	}
	return
}

func (p *Poloniex) TransferBalance(currency string, amount decimal.Decimal, from string, to string) (tb TransferBalance, err error) {
	params := url.Values{}
	params.Add("currency", currency)
	params.Add("amount", amount.StringFixed(8))
	params.Add("fromAccount", from)
	params.Add("toAccount", to)
	fmt.Printf("%+v", params)
	err = p.private("transferBalance", params, &tb)
	return
}

func (p *Poloniex) MarginAccountSummary() (mas MarginAccountSummary, err error) {
	err = p.private("returnMarginAccountSummary", nil, &mas)
	return
}

func (p *Poloniex) LoanOffer(currency string, amount decimal.Decimal, duration int, renew bool, lendingRate decimal.Decimal) (loanOffer LoanOffer, err error) {
	params := url.Values{}
	params.Add("currency", currency)
	params.Add("amount", amount.StringFixed(8))
	params.Add("lendingRate", lendingRate.Div(decimal.NewFromFloat(100.0)).StringFixed(8))
	params.Add("duration", fmt.Sprintf("%d", duration))
	r := 0
	if renew {
		r = 1
	}
	params.Add("autoRenew", fmt.Sprintf("%d", r))
	err = p.private("createLoanOffer", params, &loanOffer)
	return
}

func (p *Poloniex) CancelLoanOffer(orderNumber int64) (success bool, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	b := Base{}
	err = p.private("cancelLoanOffer", params, &b)
	success = b.Success == 1
	return
}

func (p *Poloniex) OpenLoanOffers() (openLoanOffers OpenLoanOffers, err error) {
	err = p.private("returnOpenLoanOffers", nil, &openLoanOffers)
	return
}

func (p *Poloniex) ActiveLoans() (activeLoans ActiveLoans, err error) {
	err = p.private("returnActiveLoans", nil, &activeLoans)
	provided := activeLoans.Provided
	n := []ActiveLoan{}
	for k := range provided {
		v := provided[k]
		v.Renewable = v.AutoRenew == 1
		t, err := time.Parse("2006-01-02 15:04:05", v.Date)
		if err == nil {
			v.DateTaken = t
		}
		n = append(n, v)
	}
	activeLoans.Provided = n
	return
}

func (p *Poloniex) ToggleAutoRenew(orderNumber int64) (success bool, err error) {
	params := url.Values{}
	params.Add("orderNumber", fmt.Sprintf("%d", orderNumber))
	b := Base{}
	err = p.private("toggleAutoRenew", params, &b)
	success = b.Success == 1
	return
}

// make a call to the jsonrpc api, marshal into v
func (p *Poloniex) private(method string, params url.Values, retval interface{}) error {
	if p.debug {
		defer un(trace("private: " + method))
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if params == nil {
		params = url.Values{}
	}
	params.Set("nonce", p.getNonce())
	params.Set("command", method)
	postData := params.Encode()

	req := goreq.Request{
		Method:      "POST",
		Uri:         PRIVATEURI,
		Body:        postData,
		ContentType: "application/x-www-form-urlencoded",
		Accept:      "application/json",
		Timeout:     130 * time.Second,
	}

	req.AddHeader("Sign", p.sign(postData))
	req.AddHeader("Key", p.Key)
	req.AddHeader("Content-Length", strconv.Itoa(len(postData)))

	res, err := req.Do()
	if err != nil {
		return err
	}

	defer res.Body.Close()

	s, err := res.Body.ToString()
	if err != nil {
		return err
	}

	if p.debug {
		fmt.Println(s)
	}

	// TODO: fix this shit, it's really crappy.
	if strings.HasPrefix(s, "[") {
		// poloniex only ever returns an array type when there is no real data
		// e.g. no data in a time range
		// if this ever changes then this breaks
		return nil
	}

	err = json.Unmarshal([]byte(s), retval)
	if err != nil && retval == nil {
		log.Println(err)
		return err
	}
	return err
}

// generate hmac-sha512 hash, hex encoded
func (p *Poloniex) sign(payload string) string {
	mac := hmac.New(sha512.New, []byte(p.Secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
