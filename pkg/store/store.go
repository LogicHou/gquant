package store

import (
	"errors"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	ErrNotFound = errors.New("not found")
	ErrExist    = errors.New("exist")
)

type Trader struct {
	Margin       float64
	Leverage     float64
	HistKlines   []string
	PosAmt       float64
	PosQty       int
	EntryPrice   float64
	PosSide      futures.SideType
	StopLoss     float64
	RefreshTime  map[string]int64
	Openk        float64
	LockKline    bool
	MoveStopLoss [][]float64
	MSLTrigger   float64
	Wsk
}

type Wsk struct {
	H   float64
	L   float64
	O   float64
	C   float64
	V   float64
	E   int64
	ST  int64
	Cma float64
}

type Store interface {
	Get() *Trader
	Update()
	Reset() error
}
