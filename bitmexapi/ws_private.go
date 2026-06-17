package bitmexapi

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/uws"
)

// Default recovery timeouts/limits. Copied into per-instance fields in
// NewWsPrivate; kept out of mutable package state so tests configure an instance
// rather than shared globals.
const (
	// connected but no Welcome within this window -> reconnect
	defaultWsReadyTimeout = 15 * time.Second
	// subscription sent but not confirmed within this window -> resend
	defaultWsSubscribeAckTimeout = 10 * time.Second
	// this many unconfirmed resends of a subscription -> reconnect
	defaultWsSubscribeMaxTries = 3
	// alert after this many consecutive forced reconnects (transient ones stay quiet)
	defaultWsAuthAlertThreshold = 3
	// growing re-dial delay step after repeated forced reconnects (storm guard)
	defaultWsReconnectBackoffStep = 5 * time.Second
	defaultWsReconnectBackoffMax  = time.Minute
	// watchdog tick period
	defaultWsWatchdogPeriod = 5 * time.Second
)

type WsPrivate struct {
	c              *WsClient[WsBaseResponse]
	s              *Sign
	onReady        func()
	onConnected    func()
	onDisconnected func()
	onDialError    func(error) bool
	onDialDelay    func() time.Duration
	onAuthError    func(error)
	subscriptions  *Subscriptions
	watchdog       *app.Job

	mutex        sync.Mutex
	ready        bool      // Welcome received, subscriptions sent
	connectedAt  time.Time // when the connection was established (zero - not connected)
	failures     int       // consecutive forced reconnects, for backoff/alert
	reconnecting bool      // forced reconnect in flight; dedups repeat triggers (401 on both subscribes), reset on next connect

	// recovery timings; set once in NewWsPrivate, treated as read-only afterwards
	readyTimeout         time.Duration
	subscribeAckTimeout  time.Duration
	subscribeMaxTries    int
	authAlertThreshold   int
	reconnectBackoffStep time.Duration
	reconnectBackoffMax  time.Duration
	watchdogPeriod       time.Duration
}

func NewWsPrivate(key, secret string) *WsPrivate {
	o := new(WsPrivate)
	o.c = NewWsClient[WsBaseResponse]()
	o.c.c.WithOnPreDial(o.getUrl)
	o.s = NewSign(key, secret)
	o.subscriptions = NewSubscriptions(o)
	o.subscriptions.trackPending = true // private ws uses the ACK watchdog
	o.watchdog = app.NewJob()
	o.readyTimeout = defaultWsReadyTimeout
	o.subscribeAckTimeout = defaultWsSubscribeAckTimeout
	o.subscribeMaxTries = defaultWsSubscribeMaxTries
	o.authAlertThreshold = defaultWsAuthAlertThreshold
	o.reconnectBackoffStep = defaultWsReconnectBackoffStep
	o.reconnectBackoffMax = defaultWsReconnectBackoffMax
	o.watchdogPeriod = defaultWsWatchdogPeriod
	return o
}

func (o *WsPrivate) Close() {
	o.watchdog.Stop()
	o.c.Close()
}

func (o *WsPrivate) Transport() *uws.Options {
	return o.c.Transport()
}

func (o *WsPrivate) WithLog(log *ulog.Log) *WsPrivate {
	o.c.WithLog(log)
	return o
}

func (o *WsPrivate) WithBase(base string) *WsPrivate {
	o.c.WithBase(base)
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
	// not forwarded directly to the client: our dialDelay also applies the
	// reconnect backoff on top of the caller's delay
	o.onDialDelay = f
	return o
}

func (o *WsPrivate) WithOnDialError(f func(error) bool) *WsPrivate {
	o.onDialError = f
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

// WithOnAuthError fires when the private ws stays unhealthy across several
// consecutive reconnects (rejected subscriptions, no Welcome, dial 401/403).
// A hook point for upper-layer alerting; transient self-healing stays quiet.
func (o *WsPrivate) WithOnAuthError(f func(error)) *WsPrivate {
	o.onAuthError = f
	return o
}

func (o *WsPrivate) Run() {
	o.c.WithOnConnected(func() {
		o.mutex.Lock()
		o.connectedAt = time.Now()
		o.reconnecting = false // new connection established: allow it to trigger its own reconnect
		o.mutex.Unlock()
		if o.onConnected != nil {
			o.onConnected()
		}
	})
	o.c.WithOnDisconnected(func() {
		o.mutex.Lock()
		o.ready = false
		o.connectedAt = time.Time{}
		o.mutex.Unlock()
		o.subscriptions.clearPending()
		if o.onDisconnected != nil {
			o.onDisconnected()
		}
	})
	o.c.WithOnDialDelay(o.dialDelay)
	o.c.WithOnDialError(func(err error) bool {
		s := strings.ToLower(err.Error())
		if strings.Contains(s, "401 unauthorized") || strings.Contains(s, "403 forbidden") {
			// invalid key / forbidden - stop (same as before); make it visible
			o.c.Log().Error("dial rejected, stop:", err)
			o.notifyAuthError(err)
			o.c.c.Cancel()
			return false
		}
		if o.onDialError != nil {
			return o.onDialError(err)
		}
		return false
	})
	o.c.WithOnResponse(o.onResponse)
	o.c.WithOnTopic(o.onTopic)
	o.watchdog.Tick(o.check, o.watchdogPeriod)
	o.c.Run()
}

// Ready - Welcome received and subscriptions have been sent.
func (o *WsPrivate) Ready() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
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
	switch {
	case r.IsWelcome():
		o.onWelcomed()
	case r.TokenExpired():
		// the exchange revoked the session - a re-dial with a fresh signature is needed
		o.reconnect(fmt.Errorf("token expired: [%d] %s", r.Status, r.Error))
	case r.AlreadySubscribed():
		log.Warning(r.Error)
		for _, topic := range r.RequestTopics() {
			o.subscriptions.confirm(topic)
		}
	case r.IsSubscription():
		r.Log(log)
		if r.Ok() && r.Subscribe != "" {
			o.subscriptions.confirm(r.Subscribe)
			o.resetFailures() // a confirmed subscription proves the connection works
		}
	case r.IsError():
		err := fmt.Errorf("ws error[%d] op[%s]: %s", r.Status, r.Request.Op, r.Error)
		if r.Status == 401 {
			// "User requested an account-locked subscription but no authorization
			// was provided" - the session came up unauthenticated despite the URL
			// auth params. Only reconnecting (re-rolling the dial) fixes it.
			o.reconnect(err)
		} else {
			// other subscription errors (e.g. 503 Max Pending subscription limit) -
			// the subscription stays pending and the watchdog will resend it
			log.Error(err)
		}
	default:
		r.Log(log)
	}
	return nil
}

func (o *WsPrivate) onWelcomed() {
	o.c.Log().Info("connected successfully")
	o.mutex.Lock()
	o.ready = true
	o.mutex.Unlock()
	o.subscriptions.subscribeAll()
	if o.onReady != nil {
		o.onReady()
	}
}

// reconnect drops the socket and re-dials. Counts toward the backoff/alert
// streak; the streak is reset by the first confirmed subscription.
func (o *WsPrivate) reconnect(reason error) {
	o.mutex.Lock()
	if o.reconnecting {
		// a forced reconnect is already in flight for the current connection
		// (e.g. 401 arrives on both the execution and position subscribes): count one
		// failure, alert once, drop the socket once. Reset on the next connect.
		o.mutex.Unlock()
		return
	}
	o.reconnecting = true
	o.ready = false
	o.failures++
	failures := o.failures
	o.mutex.Unlock()
	o.c.Log().Warning("reconnect:", reason)
	if failures >= o.authAlertThreshold {
		o.notifyAuthError(fmt.Errorf("private ws unhealthy (%d consecutive failures): %w", failures, reason))
	}
	o.subscriptions.clearPending()
	o.c.Reconnect()
}

func (o *WsPrivate) resetFailures() {
	o.mutex.Lock()
	o.failures = 0
	o.mutex.Unlock()
}

func (o *WsPrivate) notifyAuthError(err error) {
	if o.onAuthError != nil {
		o.onAuthError(err)
	}
}

// dialDelay - delay before re-dial: grows with each consecutive forced
// reconnect so we never storm the exchange when failures are persistent.
func (o *WsPrivate) dialDelay() time.Duration {
	o.mutex.Lock()
	failures := o.failures
	o.mutex.Unlock()
	var d time.Duration
	if failures > 1 {
		d = time.Duration(failures-1) * o.reconnectBackoffStep
		if d > o.reconnectBackoffMax {
			d = o.reconnectBackoffMax
		}
	}
	if o.onDialDelay != nil {
		if ext := o.onDialDelay(); ext > d {
			d = ext
		}
	}
	return d
}

// check - watchdog: kills stuck states that otherwise live silently while the
// connection is healthy at the transport level (ping/pong flow, no data).
// Connection state is taken from our own flags (under the mutex), not from
// c.Connected(): uws does not synchronize the socket field for outside reads.
func (o *WsPrivate) check() {
	o.mutex.Lock()
	ready := o.ready
	connectedAt := o.connectedAt
	o.mutex.Unlock()
	if connectedAt.IsZero() {
		return
	}
	if ready {
		retry, exceeded := o.subscriptions.stale(o.subscribeAckTimeout, o.subscribeMaxTries)
		if exceeded {
			o.reconnect(errors.New("subscriptions not confirmed"))
			return
		}
		for _, topic := range retry {
			o.c.Log().Warning("resubscribe:", topic)
			o.subscribe(topic)
		}
		return
	}
	if time.Since(connectedAt) >= o.readyTimeout {
		// connected, but Welcome never arrived
		o.reconnect(errors.New("no welcome after connect"))
	}
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
		base := fmt.Sprintf("%v?api-expires=%v&api-signature=%v&api-key=%v", o.c.c.Options.Base, expires, signature, o.s.Key)
		return base
	}
}
