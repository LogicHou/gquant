package dialect

import (
	"context"
	"fmt"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/indicator"
	"github.com/LogicHou/gquant/pkg/utils"
	gob "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

type binance struct {
	client *futures.Client
	conf   *config.Configuration
	logger *zap.Logger
}

var _ Dialect = (*binance)(nil)

func init() {
	Register("binance", &binance{})
}

func (b *binance) SetClient(conf *config.Configuration, logger *zap.Logger) {
	b.conf = conf
	b.client = gob.NewFuturesClient(conf.Account.AccessKey, conf.Account.SecretKey)
	b.client.NewSetServerTimeService().Do(context.Background())
}

func (b *binance) HistKlines() ([]*indicator.Kline, error) {
	klines, err := b.client.NewKlinesService().
		Symbol(b.conf.Trade.Symbol).
		Interval(b.conf.Trade.Interval).
		Limit(b.conf.Trade.HistKRange).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can not get binance kline service: %v", err)
	}

	ks := make([]*indicator.Kline, len(klines)-1)
	for i, v := range klines[:len(klines)-1] {
		kl := indicator.Kline{
			OpenTime:  v.OpenTime,
			CloseTime: v.CloseTime,
			Open:      utils.StrToF64(v.Open),
			High:      utils.StrToF64(v.High),
			Low:       utils.StrToF64(v.Low),
			Close:     utils.StrToF64(v.Close),
			Volume:    utils.StrToF64(v.Volume),
		}
		ks[i] = &kl
	}
	return ks, nil
}

func (b *binance) Ticker() (chan *indicator.Ticker, error) {
	ticker := make(chan *indicator.Ticker)

	go func() {
		_, _, err := futures.WsKlineServe(b.conf.Trade.Symbol, b.conf.Trade.Interval, func(event *futures.WsKlineEvent) {
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
			b.logger.Error("WsKlineServe handler error", zap.Error(err))
		})
		if err != nil {
			b.logger.Error("cannot start WsKlineServe", zap.Error(err))
		}

	}()

	return ticker, nil
}

func (b *binance) CreateMarketOrder(action string, price float64, qty float64, maxStopLoss float64) error {
	//TODO
	return nil
}

func (b *binance) ClosePosition(posAmt float64) error {
	//TODO
	return nil
}
