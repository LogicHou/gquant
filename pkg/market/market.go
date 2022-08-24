package market

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"syscall"

	"github.com/LogicHou/gquant/pkg/indicator"
	"github.com/LogicHou/gquant/pkg/market/kline"
	"github.com/LogicHou/gquant/pkg/market/ticker"
	"go.uber.org/zap"
)

type TickerPublisher interface {
	Publish(context.Context) error
	Subscribe() ticker.Subscriber
}

type KlinePublisher interface {
	Publish(context.Context) (chan struct{}, error)
	Subscribe() kline.Subscriber
}

type Strategy interface {
	UpdateTicker(*indicator.Ticker)
	UpdateKlines([]*indicator.Kline)
	OnKlineUpdate()
	OnTickerUpdate()
}

type Service struct {
	Logger          *zap.Logger
	Strategy        Strategy
	TickerPublisher TickerPublisher
	KlinePublisher  KlinePublisher
}

func (s *Service) Serv(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.Logger.Info("Running quant.", zap.String("market type", "futures"))

	klineSub := s.KlinePublisher.Subscribe()
	go func() {
		for {
			s.Strategy.UpdateKlines(<-klineSub)
			s.Strategy.OnKlineUpdate()
		}
	}()

	klineUpdateTrigger, err := s.KlinePublisher.Publish(ctx)
	if err != nil {
		return fmt.Errorf("can not get publish: %v", err)
	}

	klineUpdateTrigger <- struct{}{}

	tickerSub := s.TickerPublisher.Subscribe()
	go func() {
		for t := range tickerSub {
			s.Strategy.UpdateTicker(t)
			s.Strategy.OnTickerUpdate()
		}
	}()

	err = s.TickerPublisher.Publish(ctx)
	if err != nil {
		return fmt.Errorf("can not get publish: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			s.Logger.Info("ctx Done")
			return nil
		case <-c:
			s.Logger.Info("program is exiting...")
			return nil
		}
	}
}
