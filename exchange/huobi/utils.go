package huobi

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"integrator/exchange"
	"strings"
	"errors"
)

var DefaultSeparator = "."

func gzipDecode(b []byte) ([]byte, error) {

	r := bytes.NewReader(b)
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(zr)
}

func klineConverter(data KlineData, exchangeId int64) (exchange.Kline, error) {

	name := strings.Split(data.Ch, DefaultSeparator)
	pair := name[1]
	period := name[3]
	var base, quote int64
	for k, v := range exchange.AssetIdMap {
		i := strings.Index(pair, k)
		if i == 0 {
			base = v
		} else if i > 0 {
			quote = v
		}
	}
	tick := data.Tick

	kline := exchange.Kline{
		ExchangeId: exchangeId,
		BaseId:     base,
		QuoteId:    quote,
		Period:     exchange.FillPeriod(data.Ts, exchange.ConvertPeriod(period)),
		Open:       tick.Open,
		Close:      tick.Close,
		Low:        tick.Low,
		High:       tick.High,
		Volume:     tick.Amount,
		VolAmount:  tick.Count,
		Amount:     tick.Vol,
		Time:       data.Ts / 1000,
	}
	var err error
	if base == 0 || quote == 0 {
		err = errors.New("unknown asset pair")
	}
	return kline, err
}
