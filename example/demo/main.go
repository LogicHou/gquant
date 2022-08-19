package main

import (
	"context"
	"flag"
	"log"

	_ "github.com/LogicHou/gquant/example/demo/strategy"
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/market"
	"github.com/LogicHou/gquant/pkg/quant"
	"go.uber.org/zap"
)

var configPath = flag.String("c", "config.yaml", "config file path")

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	quantController := &quant.Controller{
		Name:   "binance-futures",
		Logger: logger,
		Market: &market.Market{
			Config: config.New("yaml", *configPath),
			Logger: logger,
		},
	}

	quantController.Run(context.Background())
}
