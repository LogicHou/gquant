package ticker

import (
	"context"
)

type Ticker struct {
	H float64
	L float64
	O float64
	C float64
	V float64
	E int64
}

func (t *Ticker) Subscribe(c context.Context) (ch chan *Ticker, err error) {

	return nil, nil
}
