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

func TestRsi(t *testing.T) {
	klines := getKlins()
	klen := len(klines)
	closing := make([]float64, klen)
	volume := make([]int64, klen)
	high := make([]float64, klen)
	low := make([]float64, klen)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
		volume[i] = int64(klines[i].Volume)
		high[i] = klines[i].High
		low[i] = klines[i].Low
	}
	_, rsi := RsiPeriod(6, closing)
	last10 := rsi[len(rsi)-10:]

	cases := []float64{
		59.29,
		51.91,
		61.63,
		60.32,
		54.43,
		41.51,
		44.81,
		54.87,
		70.51,
		67.56,
	}
	for i, cc := range cases {
		got := utils.FRound2(last10[9-i])
		if cc != got {
			t.Errorf("incorrect result; want: %f got: %f", cc, got)
		}
	}
}

func TestSma(t *testing.T) {
	klines := getKlins()
	klen := len(klines)
	closing := make([]float64, klen)

	for i := 0; i < klen; i++ {
		closing[i] = klines[i].Close
	}
	sma := Sma(20, closing)
	last10 := sma[len(sma)-10:]

	cases := []float64{
		1663.48,
		1661.36,
		1659.22,
		1657.64,
		1655.75,
		1653.49,
		1651.39,
		1648.98,
		1646.59,
		1644.46,
	}
	for i, cc := range cases {
		got := utils.FRound2(last10[9-i])
		if cc != got {
			t.Errorf("incorrect result; want: %f got: %f", cc, got)
		}
	}
}
