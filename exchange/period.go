package exchange

import (
	"time"
)

func GetPeriod(timestamp int64, num int, unit int) Period {

	period := Period{Num: num, Unit: unit}
	unitSec := EnumPeriodMap[unit]
	if timestamp > 9999999999 {
		timestamp /= 1000
	}

	var startAt, endAt, temp time.Time
	ts := time.Unix(timestamp, 0)
	interval := time.Duration(num*unitSec) * time.Second

	switch unitSec {
	case Min:
		startAt = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), 0, 0, 0, ts.Location())
		temp = startAt
		endAt = startAt.Add(Hour * time.Second)
	case Hour:
		startAt = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
		temp = startAt
		endAt = startAt.Add(Day * time.Second)
	case Day:
		startAt = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
		temp = startAt
		endAt = startAt.Add(Day * time.Second)
	case Week:
		startAt = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
		for startAt.Weekday() != time.Monday {
			startAt = startAt.AddDate(0, 0, -1)
		}
		temp = startAt
		endAt = startAt.Add(7 * Day * time.Second)
	default:
		return period
	}

	var timeSlice []time.Time
	for temp.Before(endAt) || temp.Equal(endAt) {
		timeSlice = append(timeSlice, temp)
		temp = temp.Add(interval)
	}

	for index, item := range timeSlice {
		next := timeSlice[index+1]
		if (ts.After(item) || ts.Equal(item)) && (ts.Before(next) || ts.Equal(next)) {
			period.Start = item.Unix()
			period.End = next.Unix()
			break
		}
	}
	return period
}

func FillPeriod(timestamp int64, period Period) Period {
	return GetPeriod(timestamp, period.Num, period.Unit)
}
