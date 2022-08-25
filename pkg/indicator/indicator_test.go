package indicator

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/LogicHou/gquant/pkg/config"
	"github.com/LogicHou/gquant/pkg/utils"
	gob "github.com/adshao/go-binance/v2"
)

var cfg = config.New("yaml", "../../example/demo/config.yaml")
var conf, _ = cfg.GetInConfig()
var client = gob.NewFuturesClient(conf.Account.AccessKey, conf.Account.SecretKey)

func getKlins() []*Kline {
	jsonb, _ := ioutil.ReadFile("./klines.json")
	if len(jsonb) > 0 {
		var ks []*Kline
		if err := json.Unmarshal(jsonb, &ks); err != nil {
			panic(err)
		}
		return ks
	}

	log.Println("start create json file")
	client.NewSetServerTimeService().Do(context.Background())
	klines, err := client.NewKlinesService().
		Symbol("ETHUSDT").
		Interval("1h").
		StartTime(1660665600000). // 2022-08-17 00:00:00
		EndTime(1661400000000).   // 2022-08-25 12:00:00
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	ks := make([]*Kline, len(klines))
	for i, v := range klines {
		ks[i] = &Kline{
			OpenTime:  v.OpenTime,
			CloseTime: v.CloseTime,
			Open:      utils.StrToF64(v.Open),
			High:      utils.StrToF64(v.High),
			Low:       utils.StrToF64(v.Low),
			Close:     utils.StrToF64(v.Close),
			Volume:    utils.StrToF64(v.Volume),
		}
	}
	var buffer bytes.Buffer
	jsonData, err := json.Marshal(ks)
	if err != nil {
		panic(err)
	}
	buffer.Write(jsonData)
	ioutil.WriteFile("./klines.json", buffer.Bytes(), 0644)
	return ks
}

func TestSma(t *testing.T) {
	klines := getKlins()
	klen := len(klines)
	closing := make([]float64, klen)
	// add cur close
	closing = append(closing, 1676.37)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
	}
	sma := Sma(22, closing)
	last10 := sma[len(sma)-10:]

	cases := []float64{
		1647.07,
		1647.33,
		1648.23,
		1649.76,
		1652.33,
		1655.09,
		1657.53,
		1658.99,
		1660.41,
		1662.70, // cur sma
	}

	for i, cc := range cases {
		got := utils.FRound2(last10[i])
		if cc != got {
			t.Errorf("incorrect result; want: %f got: %f", cc, got)
		}
	}
}

func TestRsi(t *testing.T) {
	klines := getKlins()
	klen := len(klines)
	closing := make([]float64, klen)
	// add cur close
	closing = append(closing, 1676.37)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
	}
	_, rsi := RsiPeriod(6, closing)
	last10 := rsi[len(rsi)-10:]

	cases := []float64{
		70.51,
		54.87,
		44.81,
		41.51,
		54.43,
		60.32,
		61.63,
		51.91,
		59.29,
		60.02, // cur rsi
	}
	for i, cc := range cases {
		got := utils.FRound2(last10[i])
		if cc != got {
			t.Errorf("incorrect result; want: %f got: %f", cc, got)
		}
	}
}

func TestKdj(t *testing.T) {
	klines := getKlins()
	klen := len(klines)
	closing := make([]float64, klen)
	high := make([]float64, klen)
	low := make([]float64, klen)
	// add cur close, high, low
	closing = append(closing, 1676.37)
	high = append(high, 1679.60)
	low = append(low, 1674.74)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
		high[i] = klines[i].High
		low[i] = klines[i].Low
	}

	k, d, _ := Kdj(9, 3, 3, high, low, closing)
	last10k := k[len(k)-10:]
	last10d := d[len(d)-10:]

	cases := []struct {
		k float64
		d float64
	}{
		{k: 72.96, d: 71.48},
		{k: 66.57, d: 69.85},
		{k: 55.58, d: 65.09},
		{k: 44.07, d: 58.08},
		{k: 42.88, d: 53.02},
		{k: 49.02, d: 51.69},
		{k: 54.36, d: 52.58},
		{k: 52.06, d: 52.41},
		{k: 57.85, d: 54.22},
		{k: 62.31, d: 56.92}, // cur kd
	}
	for i, cc := range cases {
		gotk := utils.FRound2(last10k[i])
		gotd := utils.FRound2(last10d[i])
		if cc.k != gotk {
			t.Errorf("incorrect result k; want: %f got: %f", cc.k, gotk)
		}
		if cc.d != gotd {
			t.Errorf("incorrect result d; want: %f got: %f", cc.d, gotd)
		}
	}
}
