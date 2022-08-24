package kline

import (
	"context"
	"fmt"

	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
)

type (
	Subscriber chan []*indicator.Kline
	Callback   func()
)

type Publisher struct {
	sub     Subscriber
	dialect dialect.Dialect
}

func NewPublisher(dialect *dialect.Dialect) *Publisher {
	return &Publisher{
		sub:     make(Subscriber),
		dialect: *dialect,
	}
}

func (p *Publisher) Subscribe() Subscriber {
	ch := make(Subscriber)
	p.sub = ch
	return ch
}

func (p *Publisher) Publish(c context.Context) (chan struct{}, error) {
	trigger := make(chan struct{})
	klines, err := p.dialect.HistKlines()
	if err != nil {
		return nil, fmt.Errorf("cannot get ticker chan: %v", err)
	}

	go func() {
		for {
			<-trigger
			p.sub <- klines
		}
	}()

	return trigger, nil
}
