package gquant

import (
	"fmt"

	"github.com/LogicHou/gquant/dialect"
	"github.com/LogicHou/gquant/indicator"
)

type Handler interface {
	Serve(*indicator.Ticker)
}

func ListenTicker(handler Handler, platform dialect.Platform) error {
	tickerCh, err := platform.Ticker()
	if err != nil {
		return fmt.Errorf("cannot get ticker chan: %v", err)
	}

	go func() {
		for t := range tickerCh {
			handler.Serve(t)
		}
	}()

	return nil
}
