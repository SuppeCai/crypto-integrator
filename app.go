package main

import (
	"encoding/json"
	"fmt"
	"integrator/exchange"
	"integrator/exchange/binance"
	"integrator/exchange/huobi"
	"integrator/status"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

var shutDown bool = false

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	status.InitAt = time.Now().Unix()
	exchange.InitLog()

	ex := exchange.Exchange{}
	ex.Execute()

	// huobi
	huobiEx := huobi.Huobi{BaseExchange: exchange.BaseExchange{Scheme: "wss", Host: "api.huobipro.com", Path: "/ws"}}
	conn, err := huobiEx.Init()
	if err != nil {
		exchange.LogErr.Error("huobi init error:" + err.Error())
	}
	defer conn.Close()
	huobiEx.Execute()

	// binance
	binanceEx := binance.Binance{BaseExchange: exchange.BaseExchange{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/stream?streams="}}
	conn, err = binanceEx.Init()
	if err != nil {
		exchange.LogErr.Error("binance init error:" + err.Error())
	}
	defer conn.Close()
	binanceEx.Execute()

	// http
	status.StartAt = time.Now().Unix()
	exchange.LogSys.Info("Integrator started!")
	go StartHttpServer()

	// data checker
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		now := time.Now().Unix()
		interval := now - status.LatestDataAt
		huobiInterval := now - huobiEx.LatestDataAt
		binanceInterval := now - binanceEx.LatestDataAt
		if huobiInterval >= 15 || shutDown {
			nf := time.Unix(now, 0)
			exchange.LogErr.Error("Lost huobi data error: Application return at:" + nf.Format("2006-01-02 15:04:05"))
			return
		}
		if binanceInterval >= 15 || shutDown {
			nf := time.Unix(now, 0)
			exchange.LogErr.Error("Lost binance data error: Application return at:" + nf.Format("2006-01-02 15:04:05"))
			return
		}
		if interval >= 15 || shutDown {
			nf := time.Unix(now, 0)
			exchange.LogErr.Error("Lost data error: Application return at:" + nf.Format("2006-01-02 15:04:05"))
			return
		}
	}
}

func StartHttpServer() {
	for {
		http.HandleFunc("/status", StatusHandler)
		http.ListenAndServe(":9090", nil)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	it := time.Unix(status.InitAt, 0)
	st := time.Unix(status.StartAt, 0)
	lt := time.Unix(status.LatestDataAt, 0)
	s := status.Status{}
	s.InitAt = it.Format("2006-01-02 15:04:05")
	s.StartAt = st.Format("2006-01-02 15:04:05")
	s.InitCost = strconv.FormatInt(status.StartAt-status.InitAt, 10) + "s"
	s.LatestDataAt = lt.Format("2006-01-02 15:04:05")
	s.LatestData = status.LatestData

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Fprintln(w, "status error!")
	} else {
		fmt.Fprintln(w, string(b))
	}
}
