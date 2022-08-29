package strategy

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

func (s *Strategy) OnTickerUpdate(ticker *indicator.Ticker) (pass bool) {
	s.ticker = ticker
	// 刷新Kline数据集
	if (s.ticker.T - s.lastKline.CloseTime) > indicator.RefreshTime[s.Conf.Trade.Interval] {
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
		if s.pid.PosAmt != 0 {
			s.pid.PosQty += 1
		}
	}

	curMa5, curMa10 := s.curIdct(s.ticker)

	// 开仓逻辑
	if s.pid.PosAmt == 0 {
		if curMa5 < curMa10 {
			s.pid.PosSide = "buy"
		} else if curMa5 > curMa10 {
			s.pid.PosSide = "sell"
		}
		fmt.Println(curMa5, curMa10, s.pid.PosSide)

		if s.openCondition(curMa5, curMa10) {
			switch s.pid.PosSide {
			case "buy":
				s.pid.StopLoss = s.lastKline.Low
			case "sell":
				s.pid.StopLoss = s.lastKline.High
			}
			s.pid.PosAmt = 0.3
			s.pid.EntryPrice = s.ticker.C
			fmt.Printf("开单，买卖方向：%s 开仓价：%f 止损位：%f 持仓数量：%f\n", s.pid.PosSide, s.ticker.C, s.pid.StopLoss, s.pid.PosAmt)
			time.Sleep(time.Second * 3)
		}
		return true
	}

	// 止盈逻辑
	if s.pid.PosQty > 1 {
		if s.tpCondition() {
			switch s.pid.PosSide {
			case "buy":
				fmt.Println("止盈 差值：", s.ticker.C-s.pid.EntryPrice)
			case "sell":
				fmt.Println("止盈 差值：", s.pid.EntryPrice-s.ticker.C)
			}
			s.resetPid()
		}
		return true
	}

	// 止损逻辑
	if s.stCondition() {
		switch s.pid.PosSide {
		case "buy":
			fmt.Println("止损 差值：", s.ticker.C-s.pid.EntryPrice)
		case "sell":
			fmt.Println("止损 差值：", s.pid.EntryPrice-s.ticker.C)
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
	s.lastKline = s.klines[len(s.klines)-1]

	s.Logger.Sugar().Infof("KlineUpdated--> PosSide:%s PosAmt:%f PosQty:%d EntryPrice:%f Leverage:%f StopLoss:%f\n", s.pid.PosSide, s.pid.PosAmt, s.pid.PosQty, s.pid.EntryPrice, s.Conf.Trade.Leverage, s.pid.StopLoss)
}
