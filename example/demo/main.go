package main

import (
	"context"
	"flag"
	"log"

	mystrategy "github.com/LogicHou/gquant/example/demo/strategy"
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/dialect"
	"github.com/LogicHou/gquant/pkg/market"
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

	tickerpub := ticker.NewPublisher(&dialect)

	mystrategy := &mystrategy.Strategy{}

	market := &market.Service{
		Strategy:        mystrategy,
		Logger:          logger,
		TickerPublisher: tickerpub,
	}

	err = market.Serv(context.Background())
	if err != nil {
		logger.Fatal("cannot start market service", zap.Error(err))
	}
}
