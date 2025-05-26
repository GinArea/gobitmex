package bitmexapi

import (
	"fmt"
	"testing"
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
