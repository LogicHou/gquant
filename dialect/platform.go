package dialect

import (
	"fmt"

	"github.com/LogicHou/gquant/indicator"
	"github.com/adshao/go-binance/v2/futures"
)

var platformsMap = map[string]Platform{}

type Platform interface {
	SetClient(*Config)
	KlineRange() ([]*indicator.Kline, error)
	Ticker() (chan *indicator.Ticker, error)
	CreateOrder(float64, indicator.ActionType, float64, float64) error
	ClosePosition(float64, float64) error
	CreateMarketOrder(indicator.ActionType, float64, float64) error
	CloseMarketPosition(float64) error
	PostionRisk() (float64, float64, float64, indicator.ActionType, error)
	GetOpenOrder() (float64, int64, error)
	GetAccountInfo() (*futures.Account, error)
}

type Config struct {
	AccessKey  string
	SecretKey  string
	Platform   string
	Symbol     string
	Interval   string
	KlineRange int
}

func Register(name string, platform Platform) {
	platformsMap[name] = platform
}

func Get(config *Config) (platform Platform, err error) {
	err = nil
	platform, ok := platformsMap[config.Platform]
	platform.SetClient(config)
	if !ok {
		err = fmt.Errorf("Cannot get platform: %v", err)
	}
	return
}
