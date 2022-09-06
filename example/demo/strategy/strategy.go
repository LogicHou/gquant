package strategy

import (
	"fmt"
	"time"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

func New(logger *zap.Logger, conf *config.Configuration, dialect dialect.Dialect, trigger chan struct{}) (strategy *Strategy) {
	s := &Strategy{
		logger:             logger,
		conf:               conf,
		dialect:            dialect,
		klineUpdateTrigger: trigger,
		lastKline:          make([]*sKline, 2),
		pid: &pid{
			PosAmt:     0.00,
			PosQty:     0,
			EntryPrice: 0.00,
			PosSide:    "",
			StopLoss:   0.00,
		},
		tune: &tune{
			CrossOffset: conf.Tune["cross_offset"].(float64),
		},
	}
	return s
}

func (s *Strategy) OnTickerUpdate(ticker *indicator.Ticker) (pass bool) {
	s.ticker = ticker
	// 刷新Kline数据集
	if (s.ticker.T - s.lastKline[0].CloseTime) > indicator.RefreshTime[s.conf.Trade.Interval] {
		oldLastOpenTime := s.lastKline[0].OpenTime
		s.klineUpdateTrigger <- struct{}{}
		time.Sleep(time.Second * 3)
		for {
			if s.lastKline[0].OpenTime == oldLastOpenTime {
				s.logger.Info("may klines update delay")
				time.Sleep(time.Second * 3)
				s.klineUpdateTrigger <- struct{}{}
			} else {
				break
			}
		}
		if s.pid.PosAmt != 0 {
			s.pid.PosQty += 1
		}
	}

	curMa5, curMa10 := s.curIdct(s.ticker)

	// 开仓逻辑
	if s.pid.PosAmt == 0 {
		if curMa5 < curMa10 {
			s.pid.PosSide = indicator.ActionBuy
		} else if curMa5 > curMa10 {
			s.pid.PosSide = indicator.ActionSell
		}
		fmt.Println(curMa5, curMa10, s.pid.PosSide)

		if s.openCondition(curMa5, curMa10) {
			switch s.pid.PosSide {
			case indicator.ActionBuy:
				s.pid.StopLoss = s.lastKline[0].Low
			case indicator.ActionSell:
				s.pid.StopLoss = s.lastKline[0].High
			}
			s.pid.PosAmt = 0.3
			s.pid.EntryPrice = s.ticker.C
			s.logger.Sugar().Infof("OP - Action: %s  EntryPrice: %f STOPLOSS: %f PosAmt: %f", s.pid.PosSide, s.pid.EntryPrice, s.pid.StopLoss, s.pid.PosAmt)
			time.Sleep(time.Second * 3)
		}
		return true
	}

	// 止盈逻辑
	if s.pid.PosQty > 1 {
		if s.tpCondition() {
			switch s.pid.PosSide {
			case indicator.ActionBuy:
				s.logger.Sugar().Infof("TP - Action: BUY  Close: %f Ratio: %f", s.ticker.C, (s.ticker.C/s.pid.EntryPrice-1)*100)
			case indicator.ActionSell:
				s.logger.Sugar().Infof("TP - Action: SELL Close: %f Ratio: %f", s.ticker.C, (s.pid.EntryPrice/s.ticker.C-1)*100)
			}
			s.resetPid()
		}
		return true
	}

	// 止损逻辑
	if s.stCondition() {
		switch s.pid.PosSide {
		case indicator.ActionBuy:
			s.logger.Sugar().Infof("ST - Action: BUY   Close: %f Ratio: %f", s.ticker.C, (s.ticker.C/s.pid.EntryPrice-1)*100)
		case indicator.ActionSell:
			s.logger.Sugar().Infof("ST - Action: SELL  Close: %f Ratio: %f", s.ticker.C, (s.ticker.C/s.pid.EntryPrice-1)*100)
		}
		s.resetPid()
	}

	return false
}

func (s *Strategy) OnKlineUpdate(klines []*indicator.Kline) {
	idct := indicator.New(klines)
	ma5 := idct.WithSma(5)
	ma10 := idct.WithSma(10)

	var sKlines = make([]*sKline, len(klines))
	for i := 0; i < len(klines); i++ {
		sKlines[i] = &sKline{klines[i], ma5[i], ma10[i]} // kline, ma5, ma10
	}

	s.klines = sKlines
	s.lastKline[0] = s.klines[len(s.klines)-1]
	s.lastKline[1] = s.klines[len(s.klines)-2]

	// s.logger.Sugar().Infof("KlineUpdated--> PosSide:%s PosAmt:%f PosQty:%d EntryPrice:%f Leverage:%f StopLoss:%f\n", s.pid.PosSide, s.pid.PosAmt, s.pid.PosQty, s.pid.EntryPrice, s.Conf.Trade.Leverage, s.pid.StopLoss)
}
