package mystrategy

import (
	"fmt"

	cst "github.com/LogicHou/gquant/pkg/const"
)

type Strategy struct {
	Name      string
	Platform  string
	AccessKey string
	SecretKey string
	Symbol    string
}

func init() {
	cc := &Strategy{
		Name:     "demo",
		Platform: cst.BINANCE,
		// AccessKey: config.,
		// SecretKey:,
		// Symbol:,
	}
	fmt.Println(cc)
}
