package main

import (
	"flag"
	"log"

	_ "github.com/LogicHou/gquant/example/demo/strategy"
	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/server"
	"go.uber.org/zap"
)

var configPath = flag.String("c", "config.yaml", "config file path")

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	err = server.Run(&server.SrvConfig{
		Name:   "binance-futures",
		Config: config.New("yaml", *configPath),
		Logger: logger,
	})
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
