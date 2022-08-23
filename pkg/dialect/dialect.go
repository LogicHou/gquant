package dialect

import (
	"fmt"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	SetClient(*config.Configuration, *zap.Logger)
	HistKlines() ([]*indicator.Kline, error)
	Ticker() (chan *indicator.Ticker, error)
	CreateMarketOrder(action string, price float64, qty float64, maxStopLoss float64) error
	ClosePosition(posAmt float64) error
}

func Register(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func Get(conf *config.Configuration) (dialect Dialect, err error) {
	err = nil
	dialect, ok := dialectsMap[conf.Dialect]
	if !ok {
		err = fmt.Errorf("Cannot get dialect: %v", err)
	}
	return
}
