package bitmexapi

import (
	"encoding/json"
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

// {"status":401,"error":"User requested an account-locked subscription but no authorization was provided.","meta":{...},"request":{"op":"subscribe","args":["execution"]}}

type WsResponse interface {
	IsSubscription() bool
	IsWelcome() bool
	IsError() bool
	IsCommandResponse() bool
	TokenExpired() bool
	AlreadySubscribed() bool
	OperationIs(string) bool
	Ok() bool
	Log(*ulog.Log)
}

// WsRequestEcho - the request echo the exchange returns in responses to op-commands.
// Args may be either an array or a single string (see the 503 example above),
// so it is parsed lazily via RawMessage.
type WsRequestEcho struct {
	Op   string          `json:"op"`
	Args json.RawMessage `json:"args"`
}

type WsBaseResponse struct {
	Success     bool          `json:"success"`
	Subscribe   string        `json:"subscribe"`
	Unsubscribe string        `json:"unsubscribe"`
	Request     WsRequestEcho `json:"request"`

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

// IsError - an error response to an op-command or a server error notification
// (e.g. 401 "no authorization was provided" on subscribe).
func (o WsBaseResponse) IsError() bool {
	return o.Status != 0 || o.Error != ""
}

// IsCommandResponse tells control responses (welcome/ack/errors) apart from topic data.
// Error responses used to fall through to processTopic and were lost with an
// "unexpected end of JSON input" - that is exactly how the lost execution/position
// subscriptions went unnoticed (zombie connection).
func (o WsBaseResponse) IsCommandResponse() bool {
	return o.IsWelcome() || o.IsError() || o.Request.Op != ""
}

// RequestTopics returns the request-echo arguments: the topics of subscribe/unsubscribe.
func (o WsBaseResponse) RequestTopics() []string {
	var list []string
	if json.Unmarshal(o.Request.Args, &list) == nil {
		return list
	}
	var single string
	if json.Unmarshal(o.Request.Args, &single) == nil && single != "" {
		return []string{single}
	}
	return nil
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
	} else if o.IsError() {
		log.Errorf("error[%d] op[%s]: %s", o.Status, o.Request.Op, o.Error)
	} else if o.Table == "" {
		log.Errorf("unhandled response: %+v", o)
	}
}
