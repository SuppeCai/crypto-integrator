package exchange

import (
	"regexp"
	"strconv"
	"bytes"
)

type CompositeKline struct {
	K       Kline
	Periods []string
}

type Kline struct {
	ExchangeId int64   `json:"exchangeId,omitempty"`
	BaseId     int64   `json:"baseId,omitempty"`
	QuoteId    int64   `json:"quoteId,omitempty"`
	Period     Period  `json:"-"`
	Open       float64 `json:"open,omitempty"`
	Close      float64 `json:"close,omitempty"`
	Low        float64 `json:"low,omitempty"`
	High       float64 `json:"high,omitempty"`
	Time       int64   `json:"time,omitempty"`
	Volume     float64 `json:"volume,omitempty"`
	VolAmount  int64   `json:"volamount,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	IsSaved    bool    `json:"isSaved"`
	UnitNum    int     `json:"unitNum,omitempty"`
	Unit       int     `json:"unit,omitempty"`
	StartAt    int64   `json:"startAt,omitempty"`
	EndAt      int64   `json:"endAt,omitempty"`
}

func (kline *Kline) GenerateKey() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(kline.ExchangeId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.FormatInt(kline.BaseId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.FormatInt(kline.QuoteId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.Itoa(kline.Period.Unit))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.Itoa(kline.Period.Num))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.FormatInt(kline.Period.Start, 10))
	return buffer.String()
}

func (kline *Kline) GeneratePrefix() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(kline.ExchangeId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.FormatInt(kline.BaseId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.FormatInt(kline.QuoteId, 10))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.Itoa(kline.Period.Unit))
	buffer.WriteString(Separator)
	buffer.WriteString(strconv.Itoa(kline.Period.Num))
	buffer.WriteString(Separator)
	buffer.WriteString("*")
	return buffer.String()
}

func (kline *Kline) PeriodInSecond() int {
	return kline.Period.Num * EnumPeriodMap[kline.Period.Unit]
}

type Klines []Kline

func (k Klines) Len() int {
	return len(k)
}
func (k Klines) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}
func (k Klines) Less(i, j int) bool {
	if k[j].Period.Start != 0 && k[i].Period.Start != 0 {
		return k[j].Period.Start < k[i].Period.Start
	} else {
		return k[j].StartAt < k[i].EndAt
	}
}

type Period struct {
	Num   int   `json:"num,omitempty"`
	Unit  int   `json:"unit,omitempty"`
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

func ConvertPeriod(p string) Period {

	numReg := regexp.MustCompile(`\d{1,3}`)
	unitReg := regexp.MustCompile(`[a-zA-z]{1,4}`)
	num, _ := strconv.Atoi(numReg.FindString(p))
	unit := unitReg.FindString(p)
	if num >= 60 && unit == "min" {
		num = 1
		unit = "hour"
	}
	return Period{Num: num, Unit: PeriodEnumMap[PeriodMap[unit]]}
}

func PeriodToString(p Period) string {

	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(p.Num))
	buffer.WriteString(PeriodNumMap[EnumPeriodMap[p.Unit]])
	return buffer.String()
}
