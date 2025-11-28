package bitmexapi

import (
	"errors"
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

type WsWalletSlice []Wallet

func (o WsWalletSlice) GetMarket() (market string) {
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

type SavedAddresses struct {
	Id                  int64
	UserId              int64
	Created             time.Time
	Currency            string
	Address             string
	Name                string
	Note                string
	SkipConfirm         bool
	SkipConfirmVerified bool
	Skip2FA             bool
	Skip2FAVerified     bool
	Network             string
	Memo                string
	CooldownExpires     time.Time
	Verified            bool
}

type GetAddress struct {
}

func (o *Client) GetYourAddresses() Response[[]SavedAddresses] {
	return GetAddress{}.Do(o)
}

func (o GetAddress) Do(c *Client) Response[[]SavedAddresses] {
	return Get(c, "v1/address", o, identity[[]SavedAddresses])
}

type GetDepositAddress struct {
	Currency string `json:"currency"`
	Network  string `json:"network"`
}

func (o *Client) GetDepositAddress(currency string, network string) Response[string] {

	if currency == "" || network == "" {
		return Response[string]{Error: errors.New("currency and network are required")}
	}

	return GetDepositAddress{
		Currency: currency,
		Network:  network,
	}.Do(o)
}

func (o GetDepositAddress) Do(c *Client) Response[string] {
	return Get(c, "v1/user/depositAddress", o, identity[string])
}
