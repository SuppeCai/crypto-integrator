package binance

type KlineData struct {
	Stream string `json:"stream,omitempty"`
	Data   struct {
		EventType string  `json:"e,omitempty"`
		EventTime float64 `json:"E,omitempty"`
		Symbol string  `json:"s,omitempty"`
		Kline struct {
			StartTime float64 `json:"t,omitempty"`
			CloseTime float64 `json:"T,omitempty"`
			Symbol string  `json:"s,omitempty"`
			Interval string  `json:"i,omitempty"`
			FirstTrade float64 `json:"f,omitempty"`
			LastTrade float64 `json:"L,omitempty"`
			Open string  `json:"o,omitempty"`
			Close string  `json:"c,omitempty"`
			High string  `json:"h,omitempty"`
			Low string  `json:"l,omitempty"`
			BaseVolume string  `json:"v,omitempty"`
			QuoteVolume string  `json:"q,omitempty"`
			TradeNumber float64 `json:"n,omitempty"`
			IsClosed bool    `json:"x,omitempty"`
		} `json:"k,omitempty"`
	} `json:"data,omitempty"`
}
