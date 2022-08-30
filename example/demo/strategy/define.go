package strategy

import (
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

type Strategy struct {
	Conf               *config.Configuration
	Logger             *zap.Logger
	ticker             *indicator.Ticker
	klines             []*sKline
	lastKline          *sKline
	klineUpdateTrigger chan struct{}
	pid                *pid
}

type pid struct {
	PosAmt     float64
	PosQty     int
	EntryPrice float64
	PosSide    string
	StopLoss   float64
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