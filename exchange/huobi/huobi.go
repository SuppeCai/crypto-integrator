package huobi

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"integrator/exchange"
	"net/url"
	"time"
)

var heartBeatChan = make(chan HeartBeat)
var msgTypes = []interface{}{KlineData{}, SubKlineResult{}, HeartBeat{}}

const HeartBeatCacheKey = "heartbeat:1"

type Huobi struct {
	configFile string
	config     Config
	subList    []string
	exchange.BaseExchange
	Conn *websocket.Conn
}

func (huobi *Huobi) Execute() {
	go huobi.Receive()
	go huobi.Heartbeat()
	go huobi.ConnectionChecker()
}

func (huobi *Huobi) Init() (*websocket.Conn, error) {

	if huobi.Conn != nil {
		huobi.Conn.Close()
	}

	u := url.URL{Scheme: huobi.Scheme, Host: huobi.Host, Path: huobi.Path}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		exchange.LogErr.Error("huobi init error:" + err.Error())
		return conn, err
	}
	huobi.Conn = conn

	err = huobi.initConfig()
	if err != nil {
		exchange.LogErr.Error("huobi config error:" + err.Error())
	}

	err = huobi.subscribe()
	if err != nil {
		exchange.LogErr.Error("huobi subscribe error:" + err.Error())
	}

	return conn, err
}

func (huobi *Huobi) Send(msg interface{}) {

	m, err := json.Marshal(msg)
	if err != nil {
		exchange.LogErr.Error("huobi send error:" + err.Error())
		return
	}

	huobi.Conn.WriteMessage(websocket.TextMessage, m)
	exchange.LogSys.Info("huobi send msg:" + string(m))
}

func (huobi *Huobi) Receive() {
	for {
		if huobi.Conn == nil {
			exchange.LogErr.Error("huobi conn nil error")
			return
		}

		_, message, err := huobi.Conn.ReadMessage()
		if err != nil {
			exchange.LogErr.Error("huobi read error:" + err.Error())
			return
		}

		msg, err := gzipDecode(message)

		huobi.Dispatch(msg)
	}
}

func (huobi *Huobi) Heartbeat() {
	for {
		req := <-heartBeatChan
		resp := HeartBeat{Pong: req.Ping}
		b, err := json.Marshal(resp)
		if err != nil {
			exchange.LogErr.Error("huobi marshal error:" + err.Error())
		} else {
			huobi.Conn.WriteMessage(websocket.TextMessage, b)
			exchange.RCache.Put(HeartBeatCacheKey, &req, exchange.Min)
		}
	}
}

func (huobi *Huobi) Dispatch(msg []byte) {

	m := string(msg)
	if CheckType(KlineData{}, m, true) {
		var data KlineData
		err := json.Unmarshal(msg, &data)
		if err != nil {
			exchange.LogErr.Error("huobi unmarshal error:" + err.Error())
			return
		}

		kline, err := klineConverter(data, exchange.ExchangeIdMap["huobi"])
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

		huobi.LatestDataAt = time.Now().Unix()
		huobi.LatestData = m
		exchange.KlineChan <- kline

		periods := huobi.config.S.ExtraPeriod[exchange.PeriodToString(kline.Period)]
		if len(periods) > 0 {
			com := exchange.CompositeKline{K: kline, Periods: periods}
			exchange.CompositeChan <- com
		}
	} else if CheckType(HeartBeat{}, m, false) {
		var beat HeartBeat
		err := json.Unmarshal(msg, &beat)
		if err != nil {
			exchange.LogErr.Error("huobi unmarshal error:" + err.Error())
			return
		}
		heartBeatChan <- beat
	}
}

func (huobi *Huobi) ConnectionChecker() {
	for {
		time.Sleep(15 * time.Second)
		beat := HeartBeat{}
		exchange.RCache.Get(HeartBeatCacheKey, &beat)
		if beat.Ping == 0 {
			_, err := huobi.Init()
			if err != nil {
				exchange.LogErr.Info("huobi reconnect init error:" + err.Error())
			} else {
				go huobi.Receive()
				exchange.LogErr.Info("huobi reconnected")
			}
		}
	}
}
