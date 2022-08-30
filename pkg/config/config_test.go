package config

import (
	"testing"
)

var cfg = New("yaml", "../../example/demo/config.yaml")

func TestGetInConfig(t *testing.T) {
	conf, _ := cfg.GetInConfig()
	conf.Tune = cfg.GetStringMap("tune")
}
