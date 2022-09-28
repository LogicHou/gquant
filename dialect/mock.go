package dialect

import (
	"context"
	"fmt"

	"github.com/LogicHou/gquant/indicator"
	"github.com/LogicHou/gquant/utils"
	gob "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type mock struct {
	client *futures.Client
	config Config
}

var _ Platform = (*mock)(nil)

func init() {
	Register("mock", &mock{})
}

func (b *mock) SetClient(config *Config) {
	b.config = *config
	b.client = gob.NewFuturesClient(b.config.AccessKey, b.config.SecretKey)
	b.client.NewSetServerTimeService().Do(context.Background())
}

func (b *mock) KlineRange() ([]*indicator.Kline, error) {
	klines, err := b.client.NewKlinesService().
		Symbol(b.config.Symbol).
		Interval(b.config.Interval).
		Limit(b.config.KlineRange).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can not get binance kline service: %v", err)
	}

	ks := make([]*indicator.Kline, len(klines)-1)
	for i, v := range klines[:len(klines)-1] {
		ks[i] = &indicator.Kline{
			OpenTime:  v.OpenTime,
			CloseTime: v.CloseTime,
			Open:      utils.StrToF64(v.Open),
			High:      utils.StrToF64(v.High),
			Low:       utils.StrToF64(v.Low),
			Close:     utils.StrToF64(v.Close),
			Volume:    utils.StrToF64(v.Volume),
		}
	}
	return ks, nil
}

func (b *mock) Ticker() (chan *indicator.Ticker, error) {
	ticker := make(chan *indicator.Ticker)

	_, _, err := futures.WsKlineServe(b.config.Symbol, b.config.Interval, func(event *futures.WsKlineEvent) {
		ticker <- &indicator.Ticker{
			O: utils.StrToF64(event.Kline.Open),
			C: utils.StrToF64(event.Kline.Close),
			H: utils.StrToF64(event.Kline.High),
			L: utils.StrToF64(event.Kline.Low),
			V: utils.StrToF64(event.Kline.Volume),
			T: event.Time,
			S: event.Kline.StartTime,
			E: event.Kline.EndTime,
		}
	}, func(err error) {
		err = fmt.Errorf("WsKlineServe handler error: %v", err)
	})
	if err != nil {
		return nil, fmt.Errorf("cannot start WsKlineServe: %v", err)
	}

	return ticker, nil
}

func (b *mock) CreateOrder(price float64, action indicator.ActionType, qty float64, stoploss float64) error {
	return nil
}

func (b *mock) ClosePosition(price float64, posAmt float64) error {
	return nil
}

func (b *mock) CreateMarketOrder(action indicator.ActionType, qty float64, stoploss float64) error {

	return nil
}

func (b *mock) CloseMarketPosition(posAmt float64) error {
	if posAmt == 0 {
		return fmt.Errorf("posAmt is zero")
	}

	return nil
}

func (b *mock) PostionRisk() (posAmt float64, entryPrice float64, leverage float64, posSide indicator.ActionType, err error) {
	posAmt = 0.00
	entryPrice = 0.00
	leverage = 0

	if posAmt > 0 {
		posSide = indicator.ActionBuy
	}
	if posAmt < 0 {
		posSide = indicator.ActionSell
	}
	return
}

func (b *mock) GetOpenOrder() (stopPrice float64, orderTime int64, err error) {
	stopPrice = 0.00
	orderTime = 0.00
	return
}

func (b *mock) GetAccountInfo() (*futures.Account, error) {
	res := &futures.Account{}
	res.TotalWalletBalance = "1000"
	res.TotalUnrealizedProfit = "100"

	return res, nil
}
