package exchange

import (
	"fmt"
	"github.com/gorilla/websocket"
	"integrator/status"
	"math"
	"sort"
	"time"
)

var AppName = "blockchain-integrator"

var DefaultCacheSize = 10
var DefaultChannelSize = 1000
var KlineChan = make(chan Kline, DefaultChannelSize)
var CompositeChan = make(chan CompositeKline, DefaultChannelSize)
var RCache = Cache{}
var Mq = MQ{}

type Exchange struct {
}

type BaseExchange struct {
	Scheme       string
	Host         string
	Path         string
	LatestDataAt int64
	LatestData   string
}

type WebSocketApi interface {
	Init() (*websocket.Conn, error)

	Send(msg interface{})

	Receive()

	Heartbeat()

	Dispatch(msg []byte)
}

func (exchange *Exchange) Execute() {
	RCache.Init()
	Mq.Init()
	go exchange.KlineHandler()
	go exchange.KlineCompositor()
}

func (exchange *Exchange) KlineHandler() {
	for {
		kline := <-KlineChan
		FillPeriodToKline(&kline)
		value := RCache.Put(kline.GenerateKey(), kline, DefaultCacheSize*kline.PeriodInSecond())
		Mq.send(kline)
		status.LatestDataAt = time.Now().Unix()
		status.LatestData = value
		fmt.Println(kline)
	}
}

func (exchange *Exchange) MQHandler() {

}

func (exchange *Exchange) KlineCompositor() {
	for {
		com := <-CompositeChan
		kline := com.K
		periods := com.Periods
		var klines []Kline
		RCache.MGetKline(kline.GeneratePrefix(), &klines)
		sort.Sort(Klines(klines))
		for _, period := range periods {
			var list, temp []Kline
			p := ConvertPeriod(period)
			p = FillPeriod(time.Now().Unix(), p)
			if len(klines) > p.Num {
				temp = klines[0:p.Num]
			} else {
				temp = klines
			}

			for _, t := range temp {
				if t.StartAt >= p.Start && t.EndAt <= p.End {
					list = append(list, t)
				}
			}

			if len(list) == 0 || list[len(list)-1].StartAt != p.Start {
				continue
			}

			new := Kline{}
			new.ExchangeId = kline.ExchangeId
			new.BaseId = kline.BaseId
			new.QuoteId = kline.QuoteId
			new.Time = kline.Time
			new.Open = list[len(list)-1].Open
			new.Close = list[0].Close
			new.Low = math.MaxFloat64
			for _, item := range list {
				if new.High < item.High {
					new.High = item.High
				}
				if new.Low > item.Low {
					new.Low = item.Low
				}
				new.Amount += item.Amount
				new.VolAmount += item.VolAmount
				new.Volume += item.Volume
			}
			new.Period = FillPeriod(new.Time, p)
			KlineChan <- new
		}
	}
}

func FillPeriodToKline(kline *Kline) {
	p := kline.Period
	kline.StartAt = p.Start
	kline.EndAt = p.End
	kline.Unit = p.Unit
	kline.UnitNum = p.Num
}

func (exchange *Exchange) Destroy() {
	Mq.close()
}
