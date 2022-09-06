package strategy

import (
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

type Strategy struct {
	conf               *config.Configuration
	tune               *tune
	logger             *zap.Logger
	dialect            dialect.Dialect
	ticker             *indicator.Ticker
	klines             []*sKline
	lastKline          []*sKline
	klineUpdateTrigger chan struct{}
	pid                *pid
}

type pid struct {
	PosAmt     float64
	PosQty     int
	EntryPrice float64
	PosSide    indicator.ActionType
	StopLoss   float64
}

type tune struct {
	CrossOffset float64
}

func (s *Strategy) SetStrategy(trigger chan struct{}) {
	s.pid = &pid{
		PosAmt:     0,
		PosQty:     0,
		EntryPrice: 0.00,
		PosSide:    "",
		StopLoss:   0.00,
	}
	s.klineUpdateTrigger = trigger
}

type sKline struct {
	*indicator.Kline
	ma5  float64
	ma10 float64
}
