package market

import (
	"context"
)

type Kline struct {
	H   float64
	L   float64
	O   float64
	C   float64
	V   float64
	E   int64
	ST  int64
	Cma float64
}

func (k *Kline) Subscribe(c context.Context) (ch chan *Kline, err error) {
	// conf, err := k.Config.GetInConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(conf)
	// fmt.Println(k.Config.GetStringMap("Account"))
	return nil, nil
}
