package strategy

import "github.com/LogicHou/gquant/pkg/indicator"

func (s *Strategy) curIdct(ticker *indicator.Ticker) (float64, float64) {
	sklen := len(s.klines)
	klines := make([]*indicator.Kline, sklen+1)
	for i := 0; i < sklen; i++ {
		klines[i] = s.klines[i].Kline
	}
	klines[sklen] = &indicator.Kline{
		OpenTime:  ticker.S,
		CloseTime: ticker.E,
		Open:      ticker.O,
		High:      ticker.H,
		Low:       ticker.L,
		Close:     ticker.C,
		Volume:    ticker.V,
	}
	idct := indicator.New(klines)
	curMa5 := idct.WithSma(5)
	curMa10 := idct.WithSma(10)
	return curMa5[len(curMa5)-1], curMa10[len(curMa10)-1]
}

const offset float64 = 1.50

func (s *Strategy) openCondition(curMa5 float64, curMa10 float64) bool {
	switch s.pid.PosSide {
	case "buy":
		if s.lastKline.ma5 < s.lastKline.ma10 && curMa5 > curMa10+offset {
			s.Logger.Info("BUY - ma5 up cross ma10 open condition triggered")
			return true
		}
	case "sell":
		if s.lastKline.ma5 > s.lastKline.ma10 && curMa5 < curMa10-offset {
			s.Logger.Info("SELL - ma5 down cross ma10 open condition triggered")
			return true
		}
	}

	return false
}

func (s *Strategy) tpCondition() bool {
	switch s.pid.PosSide {
	case "buy":
		if (s.ticker.C/s.pid.EntryPrice)-1 > 0.05 {
			s.Logger.Info("BUY - TP > 5% TP condition triggered")
			return true
		}
	case "sell":
		if (s.pid.EntryPrice/s.ticker.C)-1 > 0.05 {
			s.Logger.Info("SELL - TP > 5% TP condition triggered")
			return true
		}
	}

	return false
}

func (s *Strategy) stCondition() bool {
	switch s.pid.PosSide {
	case "buy":
		if s.ticker.C < s.pid.StopLoss {
			s.Logger.Sugar().Infof("BUY - ST condition triggered - PosSide:%s ticker.C:%f StopLoss:%f\n", s.pid.PosSide, s.ticker.C, s.pid.StopLoss)
			return true
		}
	case "sell":
		if s.ticker.C > s.pid.StopLoss {
			s.Logger.Sugar().Infof("SELL - ST condition triggered - PosSide:%s ticker.C:%f StopLoss:%f\n", s.pid.PosSide, s.ticker.C, s.pid.StopLoss)
			return true
		}
	}

	return false
}

func (s *Strategy) resetPid() {
	s.pid.EntryPrice = 0
	s.pid.PosAmt = 0
	s.pid.StopLoss = 0
	s.pid.PosQty = 0
}
