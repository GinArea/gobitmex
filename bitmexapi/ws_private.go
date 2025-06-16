package bitmexapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

type WsPrivate struct {
	c              *WsClient[WsBaseResponse]
	s              *Sign
	ready          bool
	onReady        func()
	onConnected    func()
	onDisconnected func()
	subscriptions  *Subscriptions
}

func NewWsPrivate(key, secret string) *WsPrivate {
	o := new(WsPrivate)
	o.c = NewWsClient[WsBaseResponse]()
	o.c.c.WithOnPreDial(o.getUrl)
	o.s = NewSign(key, secret)
	o.subscriptions = NewSubscriptions(o)
	o.onReady = func() {
		o.subscriptions.subscribeAll()
	}
	return o
}

func (o *WsPrivate) Close() {
	o.c.Close()
}

func (o *WsPrivate) Transport() *uws.Options {
	return o.c.Transport()
}

func (o *WsPrivate) WithLog(log *ulog.Log) *WsPrivate {
	o.c.WithLog(log)
	return o
}

func (o *WsPrivate) WithProxy(proxy string) *WsPrivate {
	o.c.WithProxy(proxy)
	return o
}

func (o *WsPrivate) WithLogRequest(enable bool) *WsPrivate {
	o.c.WithLogRequest(enable)
	return o
}

func (o *WsPrivate) WithLogResponse(enable bool) *WsPrivate {
	o.c.WithLogResponse(enable)
	return o
}

func (o *WsPrivate) WithOnDialDelay(f func() time.Duration) *WsPrivate {
	o.c.WithOnDialDelay(f)
	return o
}

func (o *WsPrivate) WithOnDialError(f func(error) bool) *WsPrivate {
	o.c.WithOnDialError(f)
	return o
}

func (o *WsPrivate) WithOnReady(f func()) *WsPrivate {
	o.onReady = f
	return o
}

func (o *WsPrivate) WithOnConnected(f func()) *WsPrivate {
	o.onConnected = f
	return o
}

func (o *WsPrivate) WithOnDisconnected(f func()) *WsPrivate {
	o.onDisconnected = f
	return o
}

func (o *WsPrivate) Run() {
	o.c.WithOnConnected(func() {
		if o.onConnected != nil {
			o.onConnected()
		}
	})
	o.c.WithOnDisconnected(func() {
		o.ready = false
		if o.onDisconnected != nil {
			o.onDisconnected()
		}
	})
	o.c.WithOnResponse(o.onResponse)
	o.c.WithOnTopic(o.onTopic)
	o.c.Run()
}

func (o *WsPrivate) Ready() bool {
	return o.ready
}

func (o *WsPrivate) subscribe(topic string) {
	o.c.Subscribe(topic)
}

func (o *WsPrivate) unsubscribe(topic string) {
	o.c.Unsubscribe(topic)
}

func (o *WsPrivate) onTopic(data []byte) error {
	return o.subscriptions.processTopic(data)
}

func (o *WsPrivate) onResponse(r WsBaseResponse) error {
	log := o.c.Log()
	if strings.HasPrefix(r.Info, "Welcome") {
		o.ready = true
		if o.onReady != nil {
			o.onReady()
		}
	}
	r.Log(log)
	return nil
}

func (o *WsPrivate) Wallet() *Executor[WsWalletSlice] {
	return NewExecutor[WsWalletSlice]("wallet", "", o.subscriptions)
}

func (o *WsPrivate) Orders() *Executor[WsOrderDetailSlice] {
	return NewExecutor[WsOrderDetailSlice]("order", "", o.subscriptions)
}

func (o *WsPrivate) Executions() *Executor[WsTradeHistorySlice] {
	return NewExecutor[WsTradeHistorySlice]("execution", "", o.subscriptions)
}

func (o *WsPrivate) Positions() *Executor[WsPositionSlice] {
	return NewExecutor[WsPositionSlice]("position", "", o.subscriptions)
}

func (o *WsPrivate) getUrl(string) string {
	if o.s == nil {
		return WebsocketUrl
	} else {
		signature, expires := o.s.GetWsSignData()
		base := fmt.Sprintf("%v?api-expires=%v&api-signature=%v&api-key=%v", WebsocketUrl, expires, signature, o.s.Key)
		ulog.Tracef("PreDial: %v", base)
		return base
	}
}
