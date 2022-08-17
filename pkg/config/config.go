package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var cfgReader *configReader

type (
	Configuration struct {
		DatabaseSettings
	}
	// 数据库配置
	DatabaseSettings struct {
		DatabaseURI  string
		DatabaseName string
		Username     string
		Password     string
	}
	// reader
	configReader struct {
		configFile string
		v          *viper.Viper
	}
)

// 获得所有配置
func GetInConfig() (configuration *Configuration, err error) {
	if err = cfgReader.v.ReadInConfig(); err != nil {
		fmt.Printf("配置文件读取失败 : %s", err)
		return nil, err
	}

	err = cfgReader.v.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("解析配置文件到结构体失败 : %s", err)
		return nil, err
	}

	return configuration, err
}

func GetStringMap(key string) map[string]interface{} {
	return cfgReader.v.GetStringMap(key)
}

// 实例化configReader
func NewConfigReader(configFile string) {
	v := viper.GetViper()
	v.SetConfigType("yaml")
	v.SetConfigFile(configFile)
	cfgReader = &configReader{
		configFile: configFile,
		v:          v,
	}
}
