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
