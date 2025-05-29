package bitmexapi

type Category string

const (
	Spot        Category = "SPOT"
	Derivatives Category = "DERIVATIVES" // USDT-M, COIN-M
)

type Bin string

const (
	Bin1m Bin = "1m"
	Bin5m Bin = "5m"
	Bin1h Bin = "1h"
	Bin1d Bin = "1d"
)

type Side string

const (
	Buy  Side = "Buy"
	Sell Side = "Sell"
)

type TimeInForce string

const (
	FillOrKill        TimeInForce = "FillOrKill"
	ImmediateOrCancel TimeInForce = "ImmediateOrCancel"
	Day               TimeInForce = "Day"
	GoodTillCancel    TimeInForce = "GoodTillCancel"
	AtTheClose        TimeInForce = "AtTheClose"
)

type OrderType string

const (
	Limit  OrderType = "Limit"
	Market OrderType = "Market"
)
