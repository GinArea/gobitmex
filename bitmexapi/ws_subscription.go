package bitmexapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Subscriptions struct {
	c     SubscriptionClient
	mutex sync.Mutex
	funcs SubscriptionFuncs
	// trackPending enables ACK tracking (private ws watchdog only). The public
	// ws leaves it false, so its subscribe path is unchanged from before the fix.
	trackPending bool
	pending      map[string]*pendingSubscription
}

// pendingSubscription - subscription was sent, ACK from the exchange not yet received
type pendingSubscription struct {
	since time.Time
	tries int
}

func NewSubscriptions(c SubscriptionClient) *Subscriptions {
	o := new(Subscriptions)
	o.c = c
	o.funcs = make(SubscriptionFuncs)
	o.pending = make(map[string]*pendingSubscription)
	return o
}

func (o *Subscriptions) subscribe(topic string, f SubscriptionFunc) {
	if o.c.Ready() {
		if o.trackPending {
			o.markPending(topic)
		}
		o.c.subscribe(topic)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.funcs[topic] = f
}

func (o *Subscriptions) unsubscribe(topic string) {
	if o.c.Ready() {
		o.c.unsubscribe(topic)
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.funcs, topic)
	delete(o.pending, topic)
}

func (o *Subscriptions) subscribeAll() {
	o.mutex.Lock()
	topics := make([]string, 0, len(o.funcs))
	for topic := range o.funcs {
		topics = append(topics, topic)
		if o.trackPending {
			o.pending[topic] = &pendingSubscription{since: time.Now()}
		}
	}
	o.mutex.Unlock()
	for _, topic := range topics {
		o.c.subscribe(topic)
	}
}

func (o *Subscriptions) markPending(topic string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.pending[topic] = &pendingSubscription{since: time.Now()}
}

// confirm clears a subscription from the ACK-waiting set
func (o *Subscriptions) confirm(topic string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.pending, topic)
}

func (o *Subscriptions) clearPending() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.pending = make(map[string]*pendingSubscription)
}

// stale returns topics not confirmed within age, to be resent. exceeded=true means
// some topic has gone unconfirmed maxTries in a row - further retries are pointless,
// a reconnect is needed.
func (o *Subscriptions) stale(age time.Duration, maxTries int) (retry []string, exceeded bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	now := time.Now()
	for topic, p := range o.pending {
		if now.Sub(p.since) >= age {
			if p.tries >= maxTries {
				exceeded = true
				return
			}
			p.tries++
			p.since = now
			retry = append(retry, topic)
		}
	}
	return
}

func (o *Subscriptions) processTopic(data []byte) (err error) {
	var topic RawTopic
	err = json.Unmarshal(data, &topic)
	if err == nil {
		blocks := o.getFunctions(topic.Table)
		if len(blocks) == 0 {
			err = fmt.Errorf("subscriptions of topic[%s] not found", topic.Table)
		} else {
			for _, block := range blocks {
				err = block.f(topic, block.market)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (o *Subscriptions) getFunctions(table string) (blocks []SubscriptioinBlock) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for key, fn := range o.funcs {
		/*
			For orderbook subscriptions key is (for example) -> orderBookL2_25:ETHUSDT
			`table` value inside raw message is -> orderBookL2_25
			We need to find all subscribed functions:
		*/
		if strings.HasPrefix(key, table) {
			parts := strings.Split(key, ":")
			var market string
			if len(parts) > 1 {
				market = parts[1]
			}

			blocks = append(blocks, SubscriptioinBlock{
				f:      fn,
				market: market,
			})
		}
	}
	return
}

// func (o *Subscriptions) getFunc(name string) (f SubscriptionFunc) {
// 	o.mutex.Lock()
// 	defer o.mutex.Unlock()
// 	for topic, fn := range o.funcs {
// 		if strings.HasPrefix(name, topic) {
// 			f = fn
// 			break
// 		}
// 	}
// 	return
// }

type SubscriptionClient interface {
	Ready() bool
	subscribe(string)
	unsubscribe(string)
}

type SubscriptioinBlock struct {
	f      SubscriptionFunc
	market string
}

type SubscriptionFunc func(RawTopic, string) error

type SubscriptionFuncs map[string]SubscriptionFunc
