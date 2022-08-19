package market

import (
	"context"
	"fmt"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/market/ticker"
	"go.uber.org/zap"
)

type TickerSubscribe interface {
	Subscribe(ctx context.Context) (ch chan *ticker.Ticker, err error)
}

type KlineSubscribe interface {
	Subscribe(ctx context.Context) (ch chan *Kline, err error)
}

type Market struct {
	Dialect         *dialect.Dialect
	Symbol          string
	Type            string
	Config          *config.Cfg
	Logger          *zap.Logger
	TickerSubscribe TickerSubscribe
	KlineSubscribe  KlineSubscribe
}

func (m *Market) Serv(c context.Context) error {
	m.Logger.Info("Running quant.", zap.String("name", "binance-futures"))
	conf, err := m.Config.GetInConfig()
	if err != nil {
		m.Logger.Error("cannot get config", zap.Error(err))
		return err
	}
	// fmt.Println(m.Config.GetStringMap("Account"))
	fmt.Println(conf)
	dialect, err := dialect.Get(conf.Dialect)
	if err != nil {
		m.Logger.Fatal("cannot get dialect", zap.Error(err))
		return err
	}
	fmt.Println(dialect)
	return nil
}
