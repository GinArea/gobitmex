package bitmexapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type Subscriptions struct {
	c     SubscriptionClient
	mutex sync.Mutex
	funcs SubscriptionFuncs
}

func NewSubscriptions(c SubscriptionClient) *Subscriptions {
	o := new(Subscriptions)
	o.c = c
	o.funcs = make(SubscriptionFuncs)
	return o
}

func (o *Subscriptions) subscribe(topic string, f SubscriptionFunc) {
	if o.c.Ready() {
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
}

func (o *Subscriptions) subscribeAll() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for topic := range o.funcs {
		o.c.subscribe(topic)
	}
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
