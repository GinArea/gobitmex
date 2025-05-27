package bitmexapi

import (
	"fmt"
	"testing"
	"time"
)

func Test_GetInstrumentActive(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "Get Instrument Active",
			client: NewClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.GetInstrumentActive(GetInstrumentActive{})
			fmt.Printf("%v", got)
		})
	}
}

func Test_GetSingleInstrument(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		instrument GetInstrument
	}{
		{
			name:   "Get Instrument Active",
			client: NewClient(),
			instrument: GetInstrument{
				Symbol: "ADAUSD",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.GetInstrument(tt.instrument)
			fmt.Printf("%v", got)
		})
	}
}

func Test_GetCandles(t *testing.T) {

	startTime, _ := time.Parse(time.RFC3339, "2025-05-01T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2025-05-15T00:00:00.000Z")

	tests := []struct {
		name   string
		client *Client
		query  GetCandle
	}{
		{
			name:   "Get XBT candles test",
			client: NewClient(),
			query: GetCandle{
				Symbol:    "XBTUSDT",
				BinSize:   Bin1d,
				Reverse:   true,
				Partial:   true,
				Count:     10,
				StartTime: &startTime,
				EndTime:   &endTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.GetCandles(tt.query)
			fmt.Printf("candles are: %v", got)
		})
	}
}
