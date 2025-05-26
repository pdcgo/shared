package common

import "time"

func GetDate(data time.Time) time.Time {
	ts := data
	day := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	return day
}
