package mystrategy

import (
	"fmt"

	"github.com/LogicHou/gquant/pkg/indicator"
)

type Strategy struct {
	ticker *indicator.Ticker
	klines []*indicator.Kline
}

func (s *Strategy) UpdateTicker(ticker *indicator.Ticker) {
	s.ticker = ticker
}

func (s *Strategy) UpdateKlines(klines []*indicator.Kline) {
	s.klines = klines
}

func (s *Strategy) OnKlineUpdate() {
	fmt.Printf("%+v\n", s.klines[len(s.klines)-1])
}

func (s *Strategy) OnTickerUpdate() {
	fmt.Println("OnTickerUpdate:", s.ticker)
}
