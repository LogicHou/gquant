package server

import (
	"fmt"

	"github.com/LogicHou/gquant/pkg/config"
	"go.uber.org/zap"
)

type SrvConfig struct {
	Name   string
	Config *config.Cfg
	Logger *zap.Logger
}

func Run(c *SrvConfig) error {
	conf, err := c.Config.GetInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	fmt.Println(c.Config.GetStringMap("Account"))
	c.Logger.Info("server started", zap.String("name", c.Name))
	return nil
}
