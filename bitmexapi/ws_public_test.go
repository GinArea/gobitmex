package bitmexapi

import (
	"fmt"
	"testing"
	"time"
)

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
		name    string
		symbol1 string
		symbol2 string
	}{
		{
			name:    "Get candles",
			symbol1: "XBTUSDT",
			symbol2: "ETHUSDT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsClient := NewWsPublic()
			wsClient.Run()
			time.Sleep(time.Duration(time.Second * 5))
			wsClient.Candles(tt.symbol1, Bin1m).Subscribe(func(v Topic[WsCandleSlice]) {
				fmt.Printf("[%v] %v\n\n", tt.symbol1, v)
			})
			wsClient.Candles(tt.symbol2, Bin1m).Subscribe(func(v Topic[WsCandleSlice]) {
				fmt.Printf("[%v] %v\n\n", tt.symbol2, v)
			})
			time.Sleep(time.Duration(time.Second * 3600))
		})
	}
}
