package main

import (
	"os"

	"github.com/LogicHou/gquant"
	_ "github.com/LogicHou/gquant/example/demo/strategy"
)

func main() {
	configFile := os.Args[1]
	gquant.Run(configFile)
}
