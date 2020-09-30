package binance

import (
	"errors"
	"integrator/exchange"
	"strconv"
	"strings"
)

var DefaultSeparator = "."

func klineConverter(data KlineData, exchangeId int64) (exchange.Kline, error) {

	pair := strings.ToLower(data.Data.Symbol)
	var base, quote int64
	for k, v := range exchange.AssetIdMap {
		i := strings.Index(pair, k)
		if i == 0 {
			base = v
		} else if i > 0 {
			quote = v
		}
	}

	k := data.Data.Kline
	open, _ := strconv.ParseFloat(k.Open, 64)
	close, _ := strconv.ParseFloat(k.Close, 64)
	low, _ := strconv.ParseFloat(k.Low, 64)
	high, _ := strconv.ParseFloat(k.High, 64)
	Volume, _ := strconv.ParseFloat(k.BaseVolume, 64)
	time := int64(data.Data.EventTime / 1000)

	kline := exchange.Kline{
		ExchangeId: exchangeId,
		BaseId:     base,
		QuoteId:    quote,
		Period:     exchange.FillPeriod(time, exchange.ConvertPeriod(data.Data.Kline.Interval)),
		Open:       open,
		Close:      close,
		Low:        low,
		High:       high,
		Volume:     Volume,
		VolAmount:  int64(k.TradeNumber),
		Amount:     Volume * close,
		Time:       time,
	}
	var err error
	if base == 0 || quote == 0 {
		err = errors.New("unknown asset pair")
	}
	return kline, err
}
