package bitmexapi

import (
	"fmt"
	"time"

	"github.com/msw-x/moon/ujson"
)

type PlaceOrder struct {
	Symbol     string
	Side       Side
	OrderQty   float64 // in units of the instrument (i.e. contracts, for spot it is the base currency in minor currency (e.g. XBt quantity for XBT)).
	ClOrdID    string
	OrdType    OrderType
	Type       TimeInForce `json:",omitempty"`
	ReduceOnly *bool       `json:",omitempty"`
	Text       string      `json:",omitempty"`
}

type OrderDetail struct {
	OrderID          string        `json:"orderID"`
	ClOrdID          string        `json:"clOrdID"`
	Account          int64         `json:"account"`
	Symbol           string        `json:"symbol"`
	Side             string        `json:"side"`
	OrderQty         ujson.Float64 `json:"orderQty"`
	DisplayQty       ujson.Float64 `json:"displayQty"`
	StopPx           ujson.Float64 `json:"stopPx"`
	PegOffsetValue   ujson.Float64 `json:"pegOffsetValue"`
	PegPriceType     string        `json:"pegPriceType"`
	Currency         string        `json:"currency"`
	SettlCurrency    string        `json:"settlCurrency"`
	OrdType          string        `json:"ordType"`
	TimeInForce      string        `json:"timeInForce"`
	ExecInst         string        `json:"execInst"`
	OrdStatus        string        `json:"ordStatus"`
	Triggered        string        `json:"triggered"`
	WorkingIndicator bool          `json:"workingIndicator"`
	OrdRejReason     string        `json:"ordRejReason"`
	LeavesQty        ujson.Float64 `json:"leavesQty"`
	CumQty           ujson.Float64 `json:"cumQty"`
	AvgPx            ujson.Float64 `json:"avgPx"`
	Text             string        `json:"text"`
	TransactTime     time.Time     `json:"transactTime"`
	Timestamp        time.Time     `json:"timestamp"`

	ClOrdLinkID     string        `json:"clOrdLinkID"`
	Price           ujson.Float64 `json:"price"`
	ContingencyType string        `json:"contingencyType"`
}

func (o *Client) PlaceOrder(v PlaceOrder) Response[OrderDetail] {
	v.Text = GinAreaTag
	return v.Do(o)
}

func (o PlaceOrder) Do(c *Client) Response[OrderDetail] {
	return Post(c, "v2/order", o, identity[OrderDetail])
}

type GetOrder struct {
	Symbol  string
	СlOrdID string
}

type GetOrderRequest struct {
	Symbol string
	Filter string
}

func (o *Client) GetOrderDetail(v GetOrder) Response[[]OrderDetail] {
	return v.DoOrderDetail(o)
}

func (o GetOrder) DoOrderDetail(c *Client) Response[[]OrderDetail] {

	req := GetOrderRequest{
		Symbol: o.Symbol,
		Filter: fmt.Sprintf("{\"clOrdID\":\"%v\"}", o.СlOrdID),
	}

	return Get(c, "v1/order", req, identity[[]OrderDetail])
}

type TradeHistory struct {
	ExecID                string        `json:"execID"`
	OrderID               string        `json:"orderID"`
	ClOrdID               string        `json:"clOrdID"`
	ClOrdLinkID           string        `json:"clOrdLinkID"`
	Account               int64         `json:"account"`
	Symbol                string        `json:"symbol"`
	Side                  string        `json:"side"`
	LastQty               int64         `json:"lastQty"`
	LastPx                ujson.Float64 `json:"lastPx"`
	LastLiquidityInd      string        `json:"lastLiquidityInd"`
	OrderQty              int64         `json:"orderQty"`
	Price                 ujson.Float64 `json:"price"`
	DisplayQty            int64         `json:"displayQty"`
	StopPx                ujson.Float64 `json:"stopPx"`
	PegOffsetValue        ujson.Float64 `json:"pegOffsetValue"`
	PegPriceType          string        `json:"pegPriceType"`
	Currency              string        `json:"currency"`
	SettlCurrency         string        `json:"settlCurrency"`
	ExecType              string        `json:"execType"`
	OrdType               string        `json:"ordType"`
	TimeInForce           string        `json:"timeInForce"`
	ExecInst              string        `json:"execInst"`
	ContingencyType       string        `json:"contingencyType"`
	OrdStatus             string        `json:"ordStatus"`
	Triggered             string        `json:"triggered"`
	WorkingIndicator      bool          `json:"workingIndicator"`
	OrdRejReason          string        `json:"ordRejReason"`
	LeavesQty             int64         `json:"leavesQty"`
	CumQty                int64         `json:"cumQty"`
	AvgPx                 ujson.Float64 `json:"avgPx"`
	Commission            ujson.Float64 `json:"commission"`
	BrokerCommission      ujson.Float64 `json:"brokerCommission"`
	FeeType               string        `json:"feeType"`
	TradePublishIndicator string        `json:"tradePublishIndicator"`
	Text                  string        `json:"text"`
	TrdMatchID            string        `json:"trdMatchID"`
	ExecCost              int64         `json:"execCost"`
	ExecComm              ujson.Float64 `json:"execComm"`
	BrokerExecComm        ujson.Float64 `json:"brokerExecComm"`
	HomeNotional          ujson.Float64 `json:"homeNotional"`
	ForeignNotional       ujson.Float64 `json:"foreignNotional"`
	TransactTime          time.Time     `json:"transactTime"`
	Timestamp             time.Time     `json:"timestamp"`
	RealisedPnl           int64         `json:"realisedPnl"`
	TrdType               string        `json:"trdType"`
}

func (o *Client) GetTradeHistory(v GetOrder) Response[[]TradeHistory] {
	return v.DoTradeHistory(o)
}

func (o GetOrder) DoTradeHistory(c *Client) Response[[]TradeHistory] {
	req := GetOrderRequest{
		Symbol: o.Symbol,
		Filter: fmt.Sprintf("{\"clOrdID\":\"%v\"}", o.СlOrdID),
	}
	return Get(c, "v1/execution/tradeHistory", req, identity[[]TradeHistory])
}
