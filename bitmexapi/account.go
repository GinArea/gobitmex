package bitmexapi

import (
	"time"

	"github.com/msw-x/moon/ujson"
)

type GetWallet struct {
	Currency string `url:"currency,omitempty"`
}

type Wallet struct {
	Account        int
	Currency       string
	Deposited      ujson.Float64
	Withdrawn      ujson.Float64
	TransferIn     ujson.Float64
	TransferOut    ujson.Float64
	Amount         ujson.Float64
	PendingCredit  ujson.Float64
	PendingDebit   ujson.Float64
	ConfirmedDebit ujson.Float64
	Timestamp      time.Time
}

type WalletShot []Wallet

func (o WalletShot) GetMarket() (market string) {
	return
}

func (o *Client) GetWalletBalance() Response[[]Wallet] {
	request := GetWallet{
		Currency: "all",
	}
	return request.Do(o)
}

func (o GetWallet) Do(c *Client) Response[[]Wallet] {
	return Get(c, "v1/user/wallet", o, identity[[]Wallet])
}
