package kline

import (
	"context"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/indicator"
	"go.uber.org/zap"
)

type Subscriber struct {
	Config *config.Cfg
	Logger *zap.Logger
}

func (s *Subscriber) Subscribe(c context.Context) (ch <-chan *indicator.Kline, err error) {
	return nil, nil
}
