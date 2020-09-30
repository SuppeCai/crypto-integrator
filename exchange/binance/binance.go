package binance

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"integrator/exchange"
	"integrator/status"
	"time"
)

const HeartBeatCacheKey = "heartbeat:3"

type Binance struct {
	configFile string
	config     Config
	subList    []string
	exchange.BaseExchange
	Conn *websocket.Conn
}

func (binance *Binance) Execute() {
	go binance.Receive()
}

func (binance *Binance) Init() (*websocket.Conn, error) {

	if binance.Conn != nil {
		binance.Conn.Close()
	}

	err := binance.initConfig()
	if err != nil {
		exchange.LogErr.Error("binance config error:" + err.Error())
	}

	path := binance.Path
	for index, value := range binance.subList {
		path += value
		if index < len(binance.subList)-1 {
			path += "/"
		}
	}

	u := binance.Scheme + "://" + binance.Host + path
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		exchange.LogErr.Error("binance init error:" + err.Error())
		return conn, err
	}
	exchange.LogSys.Info("binance multi streams:" + u)
	binance.Conn = conn
	return conn, err
}

func (binance *Binance) Receive() {
	for {
		if binance.Conn == nil {
			exchange.LogErr.Error("binance conn nil error")
			return
		}

		_, message, err := binance.Conn.ReadMessage()
		if err != nil {
			exchange.LogErr.Error("binance read error:" + err.Error())
			return
		}
		status.LatestDataAt = time.Now().Unix()
		binance.Dispatch(message)
	}
}

func (binance *Binance) Dispatch(msg []byte) {

	var data KlineData
	err := json.Unmarshal(msg, &data)
	if err != nil {
		exchange.LogErr.Error("binance unmarshal error:" + err.Error())
		return
	}

	kline, err := klineConverter(data, exchange.ExchangeIdMap["binance"])
	if err != nil {
		return
	}

	isInvalid := (kline.Open == kline.Close && kline.Open == kline.High && kline.Open == kline.Low) ||
		(kline.Open == kline.Low && kline.Close == kline.High) ||
		(kline.Open == kline.High && kline.Close == kline.Low)

	p := exchange.GetPeriod(time.Now().Unix(), kline.Period.Num, kline.Period.Unit)
	if isInvalid || p.Start != kline.Period.Start {
		return
	}

	binance.LatestDataAt = time.Now().Unix()
	binance.LatestData = string(msg)
	exchange.KlineChan <- kline
}
