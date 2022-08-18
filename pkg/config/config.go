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
		DatabaseSettings
	}
	DatabaseSettings struct {
		DatabaseURI  string
		DatabaseName string
		Username     string
		Password     string
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