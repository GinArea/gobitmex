package bitmexapi

import (
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
