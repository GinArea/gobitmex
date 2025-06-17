package bitmexapi

import (
	"fmt"
	"strings"

	"github.com/msw-x/moon/ulog"
)

// {"success":true,"subscribe":"orderBookL2_25:XBTUSD","request":{"op":"subscribe","args":["orderBookL2_25:XBTUSD"]}}

/*
{
      "table":"orderBookL2_25",
      "keys":["symbol","id","side"],
      "types":{"id":"long","price":"float","side":"symbol","size":"long","symbol":"symbol","timestamp":"timestamp"}
      "action":"partial",
      "data":[
        {"symbol":"XBTUSD","id":17999992000,"side":"Sell","size":100,"price":80,"timestamp":"2022-02-09T11:23:06.802Z"},
        {"symbol":"XBTUSD","id":17999993000,"side":"Sell","size":20,"price":70,"timestamp":"2022-02-09T11:23:06.802Z"},
        {"symbol":"XBTUSD","id":17999994000,"side":"Sell","size":10,"price":60,"timestamp":"2022-02-09T11:23:06.802Z"},
        {"symbol":"XBTUSD","id":17999995000,"side":"Buy","size":10,"price":50,"timestamp":"2022-02-09T11:23:06.802Z"},
        {"symbol":"XBTUSD","id":17999996000,"side":"Buy","size":20,"price":40,"timestamp":"2022-02-09T11:23:06.802Z"},
        {"symbol":"XBTUSD","id":17999997000,"side":"Buy","size":100,"price":30,"timestamp":"2022-02-09T11:23:06.802Z"}
      ]
    }
*/

// {"status":503,"error":"Max Pending subscription limit reached, please try again later.","request":{"op":"subscribe","args":"orderBookL2"}}

type WsResponse interface {
	IsSubscription() bool
	IsWelcome() bool
	TokenExpired() bool
	AlreadySubscribed() bool
	OperationIs(string) bool
	Ok() bool
	Log(*ulog.Log)
}

type WsBaseResponse struct {
	Success     bool        `json:"success"`
	Subscribe   string      `json:"subscribe"`
	Unsubscribe string      `json:"unsubscribe"`
	Request     interface{} `json:"request"`

	Status int    `json:"status"`
	Error  string `json:"error"`

	Info    string `json:"info"`
	AppName string `json:"appName"`

	Table  string        `json:"table"`
	Keys   []string      `json:"keys"`
	Types  interface{}   `json:"types"`
	Action string        `json:"action"`
	Data   []interface{} `json:"data"`
}

func (o WsBaseResponse) TokenExpired() bool {
	strError := strings.ToLower(o.Error)
	return o.Status == 419 && strings.Contains(strError, "access token expired")
}

func (o WsBaseResponse) AlreadySubscribed() bool {
	strError := strings.ToLower(o.Error)
	return o.Status == 400 && strings.Contains(strError, "already subscribed")
}

func (o WsBaseResponse) IsWelcome() bool {
	return o.Info != "" && o.AppName != ""
}

func (o WsBaseResponse) IsSubscription() bool {
	return o.Subscribe != "" || o.Unsubscribe != ""
}

func (o WsBaseResponse) OperationIs(v string) bool {
	return o.Table == v
}

func (o WsBaseResponse) Ok() bool {
	return o.Success
}

func (o WsBaseResponse) Log(log *ulog.Log) {
	if o.IsSubscription() {
		log.Info(fmt.Sprintf("(un)subscribe: %v", o.Success))
	} else if o.IsWelcome() {
		log.Info("connected successfully")
	} else if o.AlreadySubscribed() {
		log.Warning(o.Error)
	} else if o.TokenExpired() {
		log.Warning(o.Error)
	} else if o.Table == "" {
		log.Errorf("unhandled response: %+v", o)
	}
}
