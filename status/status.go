package status

var InitAt int64 = 0
var StartAt int64 = 0
var LatestDataAt int64 = 0
var LatestData string = "Fetching..."

type Status struct {
	InitAt       string `json:"initAt,omitempty"`
	StartAt      string `json:"startAt,omitempty"`
	InitCost     string `json:"initCost,omitempty"`
	LatestDataAt string `json:"latestDataAt,omitempty"`
	LatestData   string `json:"latestData,omitempty"`
}
