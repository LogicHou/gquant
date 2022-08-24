package kline

import (
	"context"

	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

type (
	Subscriber chan []*indicator.Kline
	Callback   func()
)

type Publisher struct {
	sub     Subscriber
	dialect dialect.Dialect
	logger  *zap.Logger
}

func NewPublisher(dialect *dialect.Dialect, logger *zap.Logger) *Publisher {
	return &Publisher{
		sub:     make(Subscriber),
		dialect: *dialect,
		logger:  logger,
	}
}

func (p *Publisher) Subscribe() Subscriber {
	ch := make(Subscriber)
	p.sub = ch
	return ch
}

func (p *Publisher) Publish(c context.Context) chan struct{} {
	trigger := make(chan struct{})

	go func() {
		for {
			<-trigger
			klines, err := p.dialect.HistKlines()
			if err != nil {
				p.logger.Error("cannot resolve account id", zap.Error(err))
			}
			p.sub <- klines
		}
	}()

	trigger <- struct{}{}

	return trigger
}
