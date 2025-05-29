package bitmexapi

import (
	"time"

	"github.com/msw-x/moon/ujson"
)

type Position struct {
	Account              int64         `json:"account"`
	Symbol               string        `json:"symbol"`
	Currency             string        `json:"currency"`
	Underlying           string        `json:"underlying"`
	QuoteCurrency        string        `json:"quoteCurrency"`
	Commission           ujson.Float64 `json:"commission"`
	InitMarginReq        ujson.Float64 `json:"initMarginReq"`
	MaintMarginReq       ujson.Float64 `json:"maintMarginReq"`
	RiskLimit            ujson.Float64 `json:"riskLimit"`
	Leverage             ujson.Float64 `json:"leverage"`
	CrossMargin          bool          `json:"crossMargin"`
	DeleveragePercentile ujson.Float64 `json:"deleveragePercentile"`
	RebalancedPnl        ujson.Float64 `json:"rebalancedPnl"`
	PrevRealisedPnl      ujson.Float64 `json:"prevRealisedPnl"`
	PrevUnrealisedPnl    ujson.Float64 `json:"prevUnrealisedPnl"`
	OpeningQty           ujson.Float64 `json:"openingQty"`
	OpenOrderBuyQty      ujson.Float64 `json:"openOrderBuyQty"`
	OpenOrderBuyCost     ujson.Float64 `json:"openOrderBuyCost"`
	OpenOrderBuyPremium  ujson.Float64 `json:"openOrderBuyPremium"`
	OpenOrderSellQty     ujson.Float64 `json:"openOrderSellQty"`
	OpenOrderSellCost    ujson.Float64 `json:"openOrderSellCost"`
	OpenOrderSellPremium ujson.Float64 `json:"openOrderSellPremium"`
	CurrentQty           ujson.Float64 `json:"currentQty"`
	CurrentCost          ujson.Float64 `json:"currentCost"`
	CurrentComm          ujson.Float64 `json:"currentComm"`
	RealisedCost         ujson.Float64 `json:"realisedCost"`
	UnrealisedCost       ujson.Float64 `json:"unrealisedCost"`
	GrossOpenPremium     ujson.Float64 `json:"grossOpenPremium"`
	IsOpen               bool          `json:"isOpen"`
	MarkPrice            ujson.Float64 `json:"markPrice"`
	MarkValue            ujson.Float64 `json:"markValue"`
	RiskValue            ujson.Float64 `json:"riskValue"`
	HomeNotional         ujson.Float64 `json:"homeNotional"`
	ForeignNotional      ujson.Float64 `json:"foreignNotional"`
	PosState             string        `json:"posState"`
	PosCost              ujson.Float64 `json:"posCost"`
	PosCross             ujson.Float64 `json:"posCross"`
	PosComm              ujson.Float64 `json:"posComm"`
	PosLoss              ujson.Float64 `json:"posLoss"`
	PosMargin            ujson.Float64 `json:"posMargin"`
	PosMaint             ujson.Float64 `json:"posMaint"`
	InitMargin           ujson.Float64 `json:"initMargin"`
	MaintMargin          ujson.Float64 `json:"maintMargin"`
	RealisedPnl          ujson.Float64 `json:"realisedPnl"`
	UnrealisedPnl        ujson.Float64 `json:"unrealisedPnl"`
	UnrealisedPnlPcnt    ujson.Float64 `json:"unrealisedPnlPcnt"`
	UnrealisedRoePcnt    ujson.Float64 `json:"unrealisedRoePcnt"`
	AvgCostPrice         ujson.Float64 `json:"avgCostPrice"`
	AvgEntryPrice        ujson.Float64 `json:"avgEntryPrice"`
	BreakEvenPrice       ujson.Float64 `json:"breakEvenPrice"`
	MarginCallPrice      ujson.Float64 `json:"marginCallPrice"`
	LiquidationPrice     ujson.Float64 `json:"liquidationPrice"`
	BankruptPrice        ujson.Float64 `json:"bankruptPrice"`
	Timestamp            time.Time     `json:"timestamp"`
}

type GetPosition struct{}

func (o *Client) GetPositions() Response[[]Position] {
	return GetPosition{}.Do(o)
}

func (o GetPosition) Do(c *Client) Response[[]Position] {
	return Get(c, "v1/position", o, identity[[]Position])
}
