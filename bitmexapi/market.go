package bitmexapi

import (
	"time"

	"github.com/msw-x/moon/ujson"
)

// GET Instuments/active
// Get all active instruments and instruments that have expired in <24hrs
// https://www.bitmex.com/api/explorer/#!/Instrument/Instrument_getActive

type GetInstrumentActive struct{}

type GetInstrument struct {
	Symbol string
}

type Instrument struct {
	Symbol                         string
	RootSymbol                     string
	State                          string
	Typ                            string
	Listing                        time.Time
	Front                          time.Time
	Expiry                         time.Time
	Settle                         time.Time
	ListedSettle                   time.Time
	PositionCurrency               string
	Underlying                     string
	QuoteCurrency                  string
	UnderlyingSymbol               string
	Reference                      string
	ReferenceSymbol                string
	MaxOrderQty                    ujson.Float64
	MaxPrice                       ujson.Float64
	LotSize                        ujson.Float64
	TickSize                       ujson.Float64
	Multiplier                     ujson.Float64
	SettlCurrency                  string
	UnderlyingToPositionMultiplier ujson.Float64
	UnderlyingToSettleMultiplier   ujson.Float64
	QuoteToSettleMultiplier        ujson.Float64
	IsQuanto                       bool
	IsInverse                      bool
	InitMargin                     ujson.Float64
	MaintMargin                    ujson.Float64
	RiskLimit                      ujson.Float64
	RiskStep                       ujson.Float64
	Limit                          ujson.Float64
	Taxed                          bool
	Deleverage                     bool
	MakerFee                       ujson.Float64
	TakerFee                       ujson.Float64
	SettlementFee                  ujson.Float64
	FundingBaseSymbol              string
	FundingQuoteSymbol             string
	FundingPremiumSymbol           string
	FundingTimestamp               time.Time
	FundingInterval                time.Time
	FundingRate                    ujson.Float64
	IndicativeFundingRate          ujson.Float64
	PrevClosePrice                 ujson.Float64
	LimitDownPrice                 ujson.Float64
	LimitUpPrice                   ujson.Float64
	PrevTotalVolume                ujson.Float64
	TotalVolume                    ujson.Float64
	Volume                         ujson.Float64
	Volume24h                      ujson.Float64
	PrevTotalTurnover              ujson.Float64
	TotalTurnover                  ujson.Float64
	Turnover                       ujson.Float64
	Turnover24h                    ujson.Float64
	HomeNotional24h                ujson.Float64
	ForeignNotional24h             ujson.Float64
	PrevPrice24h                   ujson.Float64
	Vwap                           ujson.Float64
	HighPrice                      ujson.Float64
	LowPrice                       ujson.Float64
	LastPrice                      ujson.Float64
	LastPriceProtected             ujson.Float64
	LastTickDirection              string
	LastChangePcnt                 ujson.Float64
	BidPrice                       ujson.Float64
	MidPrice                       ujson.Float64
	AskPrice                       ujson.Float64
	ImpactBidPrice                 ujson.Float64
	ImpactMidPrice                 ujson.Float64
	ImpactAskPrice                 ujson.Float64
	HasLiquidity                   bool
	OpenInterest                   ujson.Float64
	OpenValue                      ujson.Float64
	FairMethod                     string
	FairBasisRate                  ujson.Float64
	FairBasis                      ujson.Float64
	FairPrice                      ujson.Float64
	MarkMethod                     string
	MarkPrice                      ujson.Float64
	IndicativeSettlePrice          ujson.Float64
	SettledPriceAdjustmentRate     ujson.Float64
	SettledPrice                   ujson.Float64
	InstantPnl                     bool
	Timestamp                      time.Time
	MinTick                        ujson.Float64
	FundingBaseRate                ujson.Float64
	FundingQuoteRate               ujson.Float64
	Capped                         bool
	ClosingTimestamp               time.Time
	OpeningTimestamp               time.Time

	// are in the documentation but are missing from the answer:
	// CalcInterval                   time.Time
	// PublishInterval                time.Time
	// PublishTime                    time.Time
	// RebalanceTimestamp             time.Time
	// RebalanceInterval              time.Time
}

func (c *Client) GetInstrumentActive() Response[[]Instrument] {
	return GetInstrumentActive{}.Do(c)
}

func (o GetInstrumentActive) Do(c *Client) Response[[]Instrument] {
	return GetPub(c, "v1/instrument/active", o, identity[[]Instrument])
}

func (c *Client) GetInstrument(v GetInstrument) Response[[]Instrument] {
	// returns single element array
	return v.Do(c)

}

func (o GetInstrument) Do(c *Client) Response[[]Instrument] {
	return GetPub(c, "v1/instrument", o, identity[[]Instrument])
}

type Candle struct {
	Timestamp       time.Time
	Symbol          string
	Open            ujson.Float64
	High            ujson.Float64
	Low             ujson.Float64
	Close           ujson.Float64
	Trades          ujson.Int64
	Volume          ujson.Float64
	Vwap            ujson.Float64
	LastSize        ujson.Float64
	Turnover        ujson.Float64
	HomeNotional    ujson.Float64
	ForeignNotional ujson.Float64
}

type GetCandle struct {
	Symbol    string
	BinSize   Bin
	Partial   bool       `url:",omitempty"` // default: false
	Reverse   bool       `url:",omitempty"` // default: false
	Count     int        `url:",omitempty"` // default: 100
	StartTime *time.Time `url:",omitempty"`
	EndTime   *time.Time `url:",omitempty"`
	Columns   string     `url:",omitempty"`
}

func (o *Client) GetCandles(v GetCandle) Response[[]Candle] {
	return v.Do(o)
}

func (o GetCandle) Do(c *Client) Response[[]Candle] {
	return GetPub(c, "v1/trade/bucketed", o, identity[[]Candle])
}

// https://docs.bitmex.com/api-explorer/get-wallet-assets

type AssetsConfig struct {
	Asset               string        `json:"asset"`
	Currency            string        `json:"currency"`
	MajorCurrency       string        `json:"majorCurrency"`
	Name                string        `json:"name"`
	CurrencyType        string        `json:"currencyType"`
	Scale               int           `json:"scale"`
	Enabled             bool          `json:"enabled"`
	IsMarginCurrency    bool          `json:"isMarginCurrency"`
	MinDepositAmount    ujson.Float64 `json:"minDepositAmount"`
	MinWithdrawalAmount ujson.Float64 `json:"minWithdrawalAmount"`
	MaxWithdrawalAmount ujson.Float64 `json:"maxWithdrawalAmount"`
	MemoRequired        bool          `json:"memoRequired"`
	Networks            []Network     `json:"networks"`
}

type Network struct {
	Asset             string        `json:"asset"`
	TokenAddress      string        `json:"tokenAddress"`
	DepositEnabled    bool          `json:"depositEnabled"`
	WithdrawalEnabled bool          `json:"withdrawalEnabled"`
	WithdrawalFee     ujson.Float64 `json:"withdrawalFee"`
	MinFee            ujson.Float64 `json:"minFee"`
	MaxFee            ujson.Float64 `json:"maxFee"`
}

type GetAssetsConfig struct{}

func (o *Client) GetAssetsConfig() Response[[]AssetsConfig] {
	return GetAssetsConfig{}.Do(o)
}

func (o GetAssetsConfig) Do(c *Client) Response[[]AssetsConfig] {
	return GetPub(c, "v1/wallet/assets", o, identity[[]AssetsConfig])
}

type WsOrderbook struct {
	Symbol       string        `json:"symbol"`
	Id           int           `json:"id"`
	Side         string        `json:"side"`
	Size         ujson.Float64 `json:"size"`
	Price        ujson.Float64 `json:"price"`
	Timestamp    time.Time     `json:"timestamp"`
	TransactTime time.Time     `json:"transactTime"`
}

type WsOrderbookSlice []WsOrderbook

func (o WsOrderbookSlice) GetMarket() string {
	var market string
	for _, entry := range o {
		market = entry.Symbol
		break
	}
	return market
}

type WsCandleSlice []Candle

func (o WsCandleSlice) GetMarket() string {
	var market string
	for _, entry := range o {
		market = entry.Symbol
		break
	}
	return market
}
