package indicator

type Kline struct {
	OpenTime  int64
	CloseTime int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	// TradeNum  int64
	// K         float64
	// D         float64
	// MA5       float64
	// MA10      float64
	// MA20      float64
	// MA100     float64
}

type Ticker struct {
	O float64 // Open
	C float64 // Close
	H float64 // High
	L float64 // Low
	V float64 // Volume
	T int64   // Time
	S int64   // StartTime
	E int64   // EndTime
}

var RefreshTime = map[string]int64{"4h": 14404000, "1h": 3604000, "30m": 1804000, "15m": 904000, "5m": 304000, "1m": 64000}

type Indicator struct {
	closing []float64
	high    []float64
	low     []float64
	volume  []float64
	chlAvg  []float64
}

type ActionType string

const (
	ActionBuy  ActionType = "BUY"
	ActionSell ActionType = "SELL"
)

func New(klines []*Kline) *Indicator {
	klen := len(klines)
	closing := make([]float64, klen)
	high := make([]float64, klen)
	low := make([]float64, klen)
	volume := make([]float64, klen)
	chlAvg := make([]float64, klen)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
		high[i] = klines[i].High
		low[i] = klines[i].Low
		volume[i] = float64(klines[i].Volume)
		chlAvg[i] = (klines[i].Close + klines[i].High + klines[i].Low) / 3
	}

	return &Indicator{
		closing: closing,
		high:    high,
		low:     low,
		volume:  volume,
		chlAvg:  chlAvg,
	}
}

func (i *Indicator) WithVwap(period int) []float64 {
	return Vwap(period, i.chlAvg, i.volume)
}

func (i *Indicator) WithRsi(period int) []float64 {
	_, rsi := RsiPeriod(6, i.closing)
	return rsi
}

func (i *Indicator) WithKdj(rPeriod, kPeriod, dPeriod int) ([]float64, []float64) {
	k, d, _ := Kdj(rPeriod, kPeriod, dPeriod, i.high, i.low, i.closing)
	return k, d
}

func (i *Indicator) WithSma(period int) []float64 {
	return Sma(period, i.closing)
}
