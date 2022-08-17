package gquant

import (
	"fmt"

	"github.com/LogicHou/gquant/pkg/config"
)

var AppConfig = &config.Configuration{}

func Run(cfgFile string) {
	config.NewConfigReader(cfgFile)
	conf, err := config.GetInConfig()
	AppConfig = conf
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	fmt.Println(config.GetStringMap("Account"))
}
