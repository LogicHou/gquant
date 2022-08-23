package market

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
	"github.com/LogicHou/gquant/pkg/market/ticker"
	"go.uber.org/zap"
)

type TickerPublisher interface {
	Publish(context.Context) error
	Subscribe() ticker.Subscriber
}

type KlineSubscriber interface {
	Subscribe(ctx context.Context) (ch <-chan *indicator.Kline, err error)
}

type Service struct {
	Config          *config.Cfg
	Dialect         dialect.Dialect
	Logger          *zap.Logger
	TickerPublisher TickerPublisher
	KlineSubscriber KlineSubscriber
}

func (s *Service) Serv(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	s.Logger.Info("Running quant.", zap.String("market type", "futures"))
	tickerSub := s.TickerPublisher.Subscribe()
	go func() {
		for t := range tickerSub {
			fmt.Println(t)
		}
	}()
	go func() {
		time.Sleep(time.Second * 15)
		cancel()
	}()
	err := s.TickerPublisher.Publish(ctx)
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
