package dialect

import (
	"context"
	"fmt"
	"log"
	"math"

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

func (b *binance) CreateMarketOrder(action string, price float64, qty float64, maxStopLoss float64) {
	// 取消所有挂单
	err := b.client.NewCancelAllOpenOrdersService().Symbol(b.conf.Symbol).Do(context.Background())
	if err != nil {
		b.logger.Error("cannot cancel all open order", zap.Error(err))
	}
	sideStop := futures.SideTypeBuy
	sideType := futures.SideTypeSell
	offset := +5.0
	if action == indicator.ActionBuy {
		sideType = futures.SideTypeBuy
		sideStop = futures.SideTypeSell
		offset = -5.0
	}

	// 预埋止损单 RestAPI
	order, err := b.client.NewCreateOrderService().Symbol(b.conf.Symbol).
		Side(sideStop).Type("STOP_MARKET").
		ClosePosition(true).StopPrice(utils.F64ToStr(utils.FRound2(maxStopLoss + offset))).
		Do(context.Background())
	if err != nil {
		b.logger.Error("cannot create stoploss market order", zap.Error(err))
	}
	b.logger.Info("STOP_MARKET Order", zap.Any("order", order))

	// 新建市价单
	order, err = b.client.NewCreateOrderService().Symbol(b.conf.Symbol).
		Side(sideType).Type("MARKET").
		Quantity(utils.F64ToStr(qty)).
		Do(context.Background())
	if err != nil {
		b.logger.Error("cannot create market order", zap.Error(err))
	}
	b.logger.Info("MARKET Order", zap.Any("order", order))

}

func (b *binance) ClosePosition(posAmt float64) {
	if posAmt == 0 {
		b.logger.Error("posAmt is zero")
	}
	qty := posAmt

	sideType := futures.SideTypeSell
	if posAmt < 0 {
		sideType = futures.SideTypeBuy
		qty = math.Abs(posAmt)
	}

	order, err := b.client.NewCreateOrderService().Symbol(b.conf.Symbol).
		Side(sideType).Type("MARKET").
		Quantity(utils.F64ToStr(qty)).
		Do(context.Background())
	if err != nil {
		b.logger.Error("cannot create closePosition with NewCreateOrderService", zap.Error(err))
	}

	log.Println("ClosePosition:", order)

	err = b.client.NewCancelAllOpenOrdersService().Symbol(b.conf.Symbol).Do(context.Background())
	if err != nil {
		b.logger.Error("cannot create closePosition with NewCancelAllOpenOrdersService", zap.Error(err))
	}
}

func (b *binance) PostionRisk() (posAmt float64, entryPrice float64, leverage float64, posSide string) {
	res, err := b.client.NewGetPositionRiskService().Symbol(b.conf.Trade.Symbol).Do(context.Background())
	if err != nil {
		b.logger.Error("cannot get PostionRisk", zap.Error(err))
	}
	posAmt = utils.StrToF64(res[0].PositionAmt)
	entryPrice = utils.StrToF64(res[0].EntryPrice)
	leverage = utils.StrToF64(res[0].Leverage)

	if posAmt > 0 {
		posSide = indicator.ActionBuy
	}
	if posAmt < 0 {
		posSide = indicator.ActionSell
	}
	return
}

func (b *binance) GetOpenOrder() (stopPrice float64, orderTime int64) {
	res, err := b.client.NewListOpenOrdersService().Symbol(b.conf.Symbol).Do(context.Background())
	if err != nil || len(res) == 0 {
		b.logger.Error("ListOpenOrders was zero or Err", zap.Error(err))
	}
	stopPrice = utils.StrToF64(res[0].StopPrice)
	orderTime = res[0].Time
	return
}
