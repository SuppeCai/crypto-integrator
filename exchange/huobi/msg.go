package huobi

import (
	"reflect"
	//"fmt"
	//"strings"
	"strings"
)

type HeartBeat struct {
	Ping int64 `json:"ping,omitempty"`
	Pong int64 `json:"pong,omitempty"`
}

type SubKline struct {
	Sub string `json:"sub,omitempty"`
	Id  string `json:"id,omitempty"`
}

type Request struct {
	Req  string `json:"req,omitempty"`
	Id   string `json:"id,omitempty"`
	From int64  `json:"from,omitempty"`
	To   int64  `json:"to,omitempty"`
}

type SubKlineResult struct {
	Id     string `json:"id,omitempty"`
	Subbed string `json:"subbed,omitempty"`
	Ts     string `json:"ts,omitempty"`
	Status string `json:"status,omitempty"`
}

type KlineData struct {
	Ch string `json:"ch,omitempty"`
	Ts int64 `json:"ts,omitempty"`
	Tick struct {
		Id     int64 `json:"id,omitempty"`
		Open   float64 `json:"open,omitempty"`
		Close  float64 `json:"close,omitempty"`
		Low    float64 `json:"low,omitempty"`
		High   float64 `json:"high,omitempty"`
		Amount float64 `json:"amount,omitempty"`
		Vol    float64 `json:"vol,omitempty"`
		Count  int64 `json:"count,omitempty"`
	} `json:"tick,omitempty"`
}

func CheckType(s interface{}, msg string, exact bool) bool {

	var result bool
	if exact {
		result = true
	} else {
		result = false
	}

	t := reflect.TypeOf(s)
	sum := t.NumField()
	for i := 0; i < sum; i++ {
		contains := strings.Contains(strings.ToLower(msg), strings.ToLower(t.Field(i).Name))
		if exact && !contains {
			result = false
			break
		}
		if !exact && contains {
			result = true
			break
		}
	}
	return result
}
