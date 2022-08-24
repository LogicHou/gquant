package main

import (
	"context"
	"flag"
	"log"

	mystrategy "github.com/LogicHou/gquant/example/demo/strategy"
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/market"
	"github.com/LogicHou/gquant/pkg/market/kline"
	"github.com/LogicHou/gquant/pkg/market/ticker"
	"go.uber.org/zap"
)

var configPath = flag.String("c", "config.yaml", "config file path")

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	cfg := config.New("yaml", *configPath)
	conf, err := cfg.GetInConfig()
	conf.Tuning = cfg.GetStringMap("tuning")

	dialect, err := dialect.Get(conf)
	dialect.SetClient(conf, logger)

	if err != nil {
		logger.Fatal("cannot get dialect", zap.Error(err))
	}

	tickerPub := ticker.NewPublisher(&dialect)
	klinePub := kline.NewPublisher(&dialect, logger)

	mystrategy := &mystrategy.Strategy{
		Logger: logger,
		Conf:   conf,
	}

	market := &market.Service{
		Strategy:        mystrategy,
		Logger:          logger,
		TickerPublisher: tickerPub,
		KlinePublisher:  klinePub,
	}

	err = market.Serv(context.Background())
	if err != nil {
		logger.Fatal("cannot start market service", zap.Error(err))
	}
}
