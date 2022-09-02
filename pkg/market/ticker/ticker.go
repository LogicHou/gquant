package ticker

import (
	"fmt"

	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
)

type (
	Subscriber chan *indicator.Ticker
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

func (p *Publisher) Publish() error {
	ticker, err := p.dialect.Ticker()
	if err != nil {
		return fmt.Errorf("cannot get ticker chan: %v", err)
	}

	go func() {
		for t := range ticker {
			p.sub <- t
		}
	}()

	return nil
}
