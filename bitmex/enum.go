package bitmex

type Category string

const (
	Spot        Category = "SPOT"
	Derivatives Category = "DERIVATIVES" // USDT-M, COIN-M
)
