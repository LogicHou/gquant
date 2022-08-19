package quant

import (
	"context"

	"github.com/LogicHou/gquant/pkg/market"
	"go.uber.org/zap"
)

type Controller struct {
	Name   string
	Market *market.Market
	Logger *zap.Logger
}

func (c *Controller) Run(ctx context.Context) {

	c.Logger.Info("Running quant.", zap.String("name", c.Name))

	err := c.Market.Serv(ctx)
	if err != nil {
		c.Logger.Error("cannot subscribe market", zap.Error(err))
		return
	}

	for {

	}
}
