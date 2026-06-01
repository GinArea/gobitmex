package bitmexapi

import (
	"strings"
)

type Executor[T Validatable] struct {
	table         string
	market        string
	suffix        string
	subscriptions *Subscriptions
}

type Validatable interface {
	GetMarket() string
}

func NewExecutor[T Validatable](table, market string, subscriptions *Subscriptions) *Executor[T] {
	o := new(Executor[T])
	o.table = table
	o.market = market
	o.subscriptions = subscriptions
	return o
}

func (o *Executor[T]) WithSuffix(suffix string) *Executor[T] {
	o.suffix = suffix
	return o
}

func (o *Executor[T]) topicName() string {
	topic := o.table
	if o.market != "" {
		topic += ":" + o.market
	}
	if o.suffix != "" {
		topic += ":" + o.suffix
	}
	return topic
}

func (o *Executor[T]) Subscribe(onShot func(Topic[T])) {
	topic := o.topicName()
	o.subscriptions.subscribe(topic, func(raw RawTopic, market string) error {
		topic, err := UnmarshalRawTopic[T](raw)
		if err == nil {
			currentMarket := topic.Data.GetMarket()
			if strings.EqualFold(currentMarket, market) {
				onShot(topic)
			}
		}
		return err
	})
}

func (o *Executor[T]) Unsubscribe() {
	topic := o.topicName()
	o.subscriptions.unsubscribe(topic)
}
