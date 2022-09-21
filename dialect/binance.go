package dialect

import (
	"context"
	"fmt"
	"math"

	"github.com/LogicHou/gquant/indicator"
	"github.com/LogicHou/gquant/utils"
	gob "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type binance struct {
	client *futures.Client
	config Config
}

var _ Platform = (*binance)(nil)

func init() {
	Register("binance", &binance{})
}

func (b *binance) SetClient(config *Config) {
	b.config = *config
	b.client = gob.NewFuturesClient(b.config.AccessKey, b.config.SecretKey)
	b.client.NewSetServerTimeService().Do(context.Background())
}

func (b *binance) KlineRange() ([]*indicator.Kline, error) {
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

func (b *binance) Ticker() (chan *indicator.Ticker, error) {
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

func (b *binance) CreateOrder(price float64, action indicator.ActionType, qty float64, stoploss float64) error {
	// 取消所有挂单
	err := b.client.NewCancelAllOpenOrdersService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot cancel all open order: %v", err)
	}

	sideStop := futures.SideTypeBuy
	sideType := futures.SideTypeSell
	offset := +5.0
	price = price + 0.05
	if action == indicator.ActionBuy {
		sideType = futures.SideTypeBuy
		sideStop = futures.SideTypeSell
		offset = -5.0
		price = price - 0.05
	}

	// 新建限价单
	_, err = b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideType).Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).Quantity(utils.F64ToStr(qty)).
		Price(utils.F64ToStr(price)).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create market order: %v", err)
	}

	// 新建止损单
	_, err = b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideStop).Type("STOP_MARKET").
		ClosePosition(true).StopPrice(utils.F64ToStr(utils.FRound2(stoploss + offset))).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create stoploss market order: %v", err)
	}

	return nil
}

func (b *binance) ClosePosition(price float64, posAmt float64) error {
	// TODO 完成限价止损
	if posAmt == 0 {
		return fmt.Errorf("posAmt is zero")
	}

	qty := posAmt

	sideType := futures.SideTypeSell
	if posAmt < 0 {
		sideType = futures.SideTypeBuy
		qty = math.Abs(posAmt)
	}

	_, err := b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideType).Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).Quantity(utils.F64ToStr(qty)).
		Price(utils.F64ToStr(price)).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create market order: %v", err)
	}

	err = b.client.NewCancelAllOpenOrdersService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create closePosition with NewCancelAllOpenOrdersService: %v", err)
	}

	return nil
}

func (b *binance) CreateMarketOrder(action indicator.ActionType, qty float64, stoploss float64) error {
	// 取消所有挂单
	err := b.client.NewCancelAllOpenOrdersService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot cancel all open order: %v", err)
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
	_, err = b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideStop).Type("STOP_MARKET").
		ClosePosition(true).StopPrice(utils.F64ToStr(utils.FRound2(stoploss + offset))).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create stoploss market order: %v", err)
	}

	// 新建市价单
	_, err = b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideType).Type("MARKET").
		Quantity(utils.F64ToStr(qty)).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create market order: %v", err)
	}

	return nil
}

func (b *binance) CloseMarketPosition(posAmt float64) error {
	if posAmt == 0 {
		return fmt.Errorf("posAmt is zero")
	}

	qty := posAmt

	sideType := futures.SideTypeSell
	if posAmt < 0 {
		sideType = futures.SideTypeBuy
		qty = math.Abs(posAmt)
	}

	_, err := b.client.NewCreateOrderService().Symbol(b.config.Symbol).
		Side(sideType).Type("MARKET").
		Quantity(utils.F64ToStr(qty)).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create closePosition with NewCreateOrderService: %v", err)
	}

	err = b.client.NewCancelAllOpenOrdersService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create closePosition with NewCancelAllOpenOrdersService: %v", err)
	}

	return nil
}

func (b *binance) PostionRisk() (posAmt float64, entryPrice float64, leverage float64, posSide indicator.ActionType, err error) {
	res, err := b.client.NewGetPositionRiskService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("cannot get PostionRisk: %v", err)
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

func (b *binance) GetOpenOrder() (stopPrice float64, orderTime int64, err error) {
	res, err := b.client.NewListOpenOrdersService().Symbol(b.config.Symbol).Do(context.Background())
	if err != nil || len(res) == 0 {
		err = fmt.Errorf("ListOpenOrders was zero or Err: %v", err)
	}
	stopPrice = utils.StrToF64(res[0].StopPrice)
	orderTime = res[0].Time
	return
}

func (b *binance) GetAccountInfo() (*futures.Account, error) {
	res, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
