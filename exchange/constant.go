package exchange

const (
	Min  = 60
	Hour = 60 * Min
	Day  = 24 * Hour
	Week = 7 * Day
)

var PeriodMap = map[string]int{
	"min":  Min,
	"hour": Hour,
	"day":  Day,
	"week": Week,
	"m":  Min,
	"h": Hour,
	"d":  Day,
	"w": Week,
}

var PeriodNumMap = map[int]string{
	Min:  "min",
	Hour: "hour",
	Day:  "day",
	Week: "week",
}

var PeriodEnumMap = map[int]int{
	Min:  1,
	Hour: 2,
	Day:  3,
	Week: 4,
}

var EnumPeriodMap = map[int]int{
	1: Min,
	2: Hour,
	3: Day,
	4: Week,
}

var AssetIdMap = map[string]int64{
	"btc":   1,
	"eth":   2,
	"xrp":   3,
	"bch":   4,
	"eos":   5,
	"ltc":   6,
	"ada":   7,
	"xlm":   8,
	"trx":   9,
	"iota":  10,
	"neo":   11,
	"dash":  12,
	"xmr":   13,
	"xem":   14,
	"etc":   15,
	"ven":   16,
	"omg":   17,
	"zec":   18,
	"zil":   19,
	"ht":    20,
	"edu":   21,
	"iost":  22,
	"steem": 23,
	"usdt":  24,
	"lsk":  25,
	"ont":  26,
	"zrx":  27,
	"nano":  28,
	"icx":  29,
	"xvg":  30,
	"qtum":  31,
	"bsv":  32,
	"tusd":  33,
	"bnb":  34,
	"bchabc":  4,
	"rvn":  35,
}

var ExchangeIdMap = map[string]int64{
	"huobi": 1,
	"okex":  2,
	"binance":  3,
}
//var PeriodNumMap = map[string][]int{
//	"min":  {1, 5, 15, 30},
//	"hour": {1, 2, 4, 6, 12},
//	"day":  {1},
//	"week": {1},
//}
