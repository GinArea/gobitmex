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

// ReferralCodeDetails - fee discount/rebate configuration of a referral code
type ReferralCodeDetails struct {
	// TradingFeeDiscount - trading fee discount applied to referees (fraction, e.g. 0.2 = 20%)
	TradingFeeDiscount float64 `json:",omitempty"`
	// TradingFeeDiscountDuration - duration in months of the trading fee discount
	TradingFeeDiscountDuration int `json:",omitempty"`
	// TradingFeeRebate - trading fee rebate paid back to the referrer (fraction, e.g. 0.1 = 10%)
	TradingFeeRebate float64 `json:",omitempty"`
	// TradingFeeRebateDuration - duration in months of the trading fee rebate
	TradingFeeRebateDuration int `json:",omitempty"`
}

// CreateReferralCode - request for POST /api/v1/referralCode
// https://docs.bitmex.com/api-explorer/referral-code-create-referral-code
// Note: request body is FLAT (fee fields at top level), even though the response
// nests them under "details". Nesting on the request is rejected as
// "Validation Error: data should NOT have additional properties".
type CreateReferralCode struct {
	// Code - the desired referral code string
	Code string
	// IsDefault - set true to mark this code as the account's default. Only true is accepted;
	// false is rejected by the server as an enum-validation error, so false is omitted.
	IsDefault bool `json:",omitempty"`
	// IsPrimary - set true to mark this code as the account's primary. Only true is accepted;
	// false is rejected by the server as an enum-validation error, so false is omitted.
	IsPrimary bool `json:",omitempty"`
	// TradingFeeDiscount - trading fee discount applied to referees (fraction, e.g. 0.2 = 20%); zero omits
	TradingFeeDiscount float64 `json:",omitempty"`
	// TradingFeeDiscountDuration - duration in months of the trading fee discount; zero omits
	TradingFeeDiscountDuration int `json:",omitempty"`
	// TradingFeeRebate - trading fee rebate paid back to the referrer (fraction, e.g. 0.1 = 10%); zero omits
	TradingFeeRebate float64 `json:",omitempty"`
	// TradingFeeRebateDuration - duration in months of the trading fee rebate; zero omits
	TradingFeeRebateDuration int `json:",omitempty"`
}

// ReferralCode - response for POST /api/v1/referralCode
// https://docs.bitmex.com/api-explorer/referral-code-create-referral-code
type ReferralCode struct {
	// Id - unique identifier of the referral code
	Id string `json:"id"`
	// UserId - identifier of the user who owns the referral code
	UserId int `json:"userId"`
	// Code - the referral code string
	Code string `json:"code"`
	// Details - fee-discount / rebate configuration of this code
	Details ReferralCodeDetails `json:"details"`
	// Created - timestamp when the referral code was created
	Created time.Time `json:"created"`
	// Modified - timestamp when the referral code was last modified
	Modified time.Time `json:"modified"`
	// IsDefault - whether this is the default referral code
	IsDefault bool `json:"isDefault"`
	// IsPrimary - whether this is the primary referral code
	IsPrimary bool `json:"isPrimary"`
}

func (o *Client) CreateReferralCode(v CreateReferralCode) Response[ReferralCode] {
	return v.Do(o)
}

func (o CreateReferralCode) Do(c *Client) Response[ReferralCode] {
	return Post(c, "v1/referralCode", o, identity[ReferralCode])
}

// GetAllReferralCodes - request for GET /api/v1/referralCode
// https://docs.bitmex.com/api-explorer/referral-code-get-all-codes-for-user
type GetAllReferralCodes struct{}

func (o *Client) GetAllReferralCodes() Response[[]ReferralCode] {
	return GetAllReferralCodes{}.Do(o)
}

func (o GetAllReferralCodes) Do(c *Client) Response[[]ReferralCode] {
	return Get(c, "v1/referralCode", o, identity[[]ReferralCode])
}

// DeleteReferralCode - request for DELETE /api/v1/referralCode/{id}
// https://docs.bitmex.com/api-explorer/referral-code-delete-referral-code
// Success response body: {"success":true}
// Error response body:   {"error":{"message":"Not Found","name":"HTTPError"}}
type DeleteReferralCode struct {
	// Id - unique identifier of the referral code to delete (path parameter)
	Id string
}

type deleteReferralCodeResponse struct {
	Success bool `json:"success"`
}

func (o *Client) DeleteReferralCode(id string) Response[bool] {
	return DeleteReferralCode{Id: id}.Do(o)
}

func (o DeleteReferralCode) Do(c *Client) Response[bool] {
	return Delete(c, "v1/referralCode/"+o.Id, o, func(r deleteReferralCodeResponse) (bool, error) {
		return r.Success, nil
	})
}
