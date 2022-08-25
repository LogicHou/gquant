package mystrategy

import (
	"fmt"
	"time"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

type Strategy struct {
	Conf               *config.Configuration
	Logger             *zap.Logger
	klines             []*indicator.Kline
	lastKline          *indicator.Kline
	klineUpdateTrigger chan struct{}
	pid                *pid
}

type pid struct {
	PosAmt       float64
	PosQty       int
	EntryPrice   float64
	PosSide      string
	StopLoss     float64
	Openk        float64
	MoveStopLoss [][]float64
	MSLTrigger   float64
}

func (s *Strategy) SetStrategy(trigger chan struct{}) {
	s.pid = &pid{
		PosAmt:       0,
		PosQty:       0,
		EntryPrice:   0.00,
		PosSide:      "",
		StopLoss:     0.00,
		Openk:        0.00,
		MoveStopLoss: nil,
		MSLTrigger:   0.00,
	}
	s.klineUpdateTrigger = trigger
}

func (s *Strategy) OnTickerUpdate(ticker *indicator.Ticker) (pass bool) {
	if (ticker.T - s.lastKline.CloseTime) > indicator.RefreshTime[s.Conf.Trade.Interval] {
		oldLastOpenTime := s.lastKline.OpenTime
		s.klineUpdateTrigger <- struct{}{}
		time.Sleep(time.Second * 3)
		for {
			if s.lastKline.OpenTime == oldLastOpenTime {
				s.Logger.Info("may klines update delay")
				time.Sleep(time.Second * 3)
				s.klineUpdateTrigger <- struct{}{}
			} else {
				break
			}
		}
		s.Logger.Sugar().Infof("KlineUpdated--> PosSide:%s PosAmt:%f PosQty:%d EntryPrice:%f Leverage:%f StopLoss:%f\n", s.pid.PosSide, s.pid.PosAmt, s.pid.PosQty, s.pid.EntryPrice, s.Conf.Trade.Leverage, s.pid.StopLoss)
	}

	return false
}

func (s *Strategy) OnKlineUpdate(klines []*indicator.Kline) {
	s.klines = klines
	s.lastKline = s.klines[len(s.klines)-1]
	fmt.Printf("%+v\n", s.klines[len(s.klines)-1])
}
