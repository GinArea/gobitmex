package bitmexapi

import (
	"fmt"
	"testing"
	"time"
)

func TestOrderBook(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
	}{
		{
			name:   "Orderbook",
			symbol: "WOOUSDT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wsClient := NewWsPublic()
			wsClient.Run()
			time.Sleep(time.Duration(time.Second * 5))
			wsClient.Orderbook(tt.symbol).Subscribe(func(v Topic[[]WsOrderbook]) {
				fmt.Printf("%v\n\n", v)
			})
			// wsClient.Orderbook("MELANIAUSDT").Subscribe(func(v Topic[[]WsOrderbook]) {
			// 	fmt.Printf("%v\n\n", v)
			// })
			// wsClient.Orderbook(tt.symbol).Unsubscribe()
			time.Sleep(time.Duration(time.Second * 3600))
		})
	}
}
