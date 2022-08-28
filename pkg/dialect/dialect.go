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
	CreateMarketOrder(string, float64, float64, float64)
	ClosePosition(float64)
	PostionRisk() (float64, float64, float64, string)
	GetOpenOrder() (float64, int64)
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
