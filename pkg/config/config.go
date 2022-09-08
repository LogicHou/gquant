package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Cfg struct {
		v *viper.Viper
	}
	Configuration struct {
		Account
		Trade
		Tune map[string]interface{}
	}
	Account struct {
		AccessKey  string
		SecretKey  string
		TgBotToken string
	}
	Trade struct {
		Dialect     string
		Symbol      string
		Interval    string
		Leverage    float64
		HistKRange  int
		Margin      float64
		MarginRatio float64
		MarginLimit float64
	}
)

func New(cfgType string, cfgFile string) *Cfg {
	v := viper.GetViper()
	v.SetConfigType(cfgType)
	v.SetConfigFile(cfgFile)
	return &Cfg{
		v: v,
	}
}

func (c *Cfg) GetInConfig() (configuration *Configuration, err error) {
	if err = c.v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot open config file: %v", err)
	}

	err = c.v.Unmarshal(&configuration)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config to struct : %v", err)
	}

	return configuration, err
}

func (c *Cfg) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}
