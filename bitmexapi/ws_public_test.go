package bitmexapi

import (
	"fmt"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "connection testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsClient := NewWsPublic()
			wsClient.Run()
			time.Sleep(time.Duration(time.Second * 3600))
		})
	}
}

func TestOrderBook(t *testing.T) {
	tests := []struct {
		name    string
		symbol1 string
		symbol2 string
	}{
		{
			name:    "Orderbook",
			symbol1: "ETHUSDT",
			symbol2: "MELANIAUSDT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsClient := NewWsPublic()
			wsClient.Run()
			time.Sleep(time.Duration(time.Second * 5))
			wsClient.Orderbook(tt.symbol1).Subscribe(func(v Topic[WsOrderbookSlice]) {
				fmt.Printf("[%v] %v\n\n", tt.symbol1, v)
			})
			wsClient.Orderbook(tt.symbol2).Subscribe(func(v Topic[WsOrderbookSlice]) {
				fmt.Printf("[%v] %v\n\n", tt.symbol2, v)
			})
			time.Sleep(time.Duration(time.Second * 10))
			wsClient.Orderbook(tt.symbol1).Unsubscribe()
			time.Sleep(time.Duration(time.Second * 3600))
		})
	}
}

func TestCandles(t *testing.T) {
	tests := []struct {
		name      string
		symbol1   string
		interval1 Bin
		symbol2   string
		interval2 Bin
	}{
		{
			name:      "Get candles",
			symbol1:   "XBTUSDT",
			symbol2:   "ETHUSDT",
			interval1: Bin1m,
			interval2: Bin5m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsClient := NewWsPublic()
			wsClient.Run()
			time.Sleep(time.Duration(time.Second * 5))
			wsClient.Candles(tt.symbol1, tt.interval1).Subscribe(func(v Topic[WsCandleSlice]) {
				s := v.Data[0]
				fmt.Printf("[%v] action: %v\ntime: %v\nOpen: %v\nHigh: %v\nLow: %v\nClose: %v\nTrades: %v\nVwap: %v\nVolume: %v\nLastSize: %v\nTurnover: %v\nHomeNotional: %v\nForeignNotional: %v\n--------------\n", tt.symbol1, v.Action, s.Timestamp, s.Open, s.High, s.Low, s.Close, s.Trades, s.Vwap, s.Volume, s.LastSize, s.Turnover, s.HomeNotional, s.ForeignNotional)
			})
			wsClient.Candles(tt.symbol2, tt.interval2).Subscribe(func(v Topic[WsCandleSlice]) {
				s := v.Data[0]
				fmt.Printf("[%v] action: %v\ntime: %v\nOpen: %v\nHigh: %v\nLow: %v\nClose: %v\nTrades: %v\nVwap: %v\nVolume: %v\nLastSize: %v\nTurnover: %v\nHomeNotional: %v\nForeignNotional: %v\n--------------\n", tt.symbol2, v.Action, s.Timestamp, s.Open, s.High, s.Low, s.Close, s.Trades, s.Vwap, s.Volume, s.LastSize, s.Turnover, s.HomeNotional, s.ForeignNotional)
			})
			time.Sleep(time.Duration(time.Second * 3600))
		})
	}
}
