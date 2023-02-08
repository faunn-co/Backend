package utils

import "time"

const (
	//TimeZone constant timezone
	TimeZone = "Asia/Singapore"
	WEEK     = 604800
)

var (
	tz, _ = time.LoadLocation(TimeZone)
)

// UnixToUTC Converts current unix time to UTC time object
func UnixToUTC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Local().UTC()
}

// WeekStartEndDate Returns the start and end day of the current week in SGT unix time
func WeekStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)

	startOffset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	startResult := date.Add(time.Duration(startOffset*24) * time.Hour)
	endResult := startResult.Add(time.Duration(6*24) * time.Hour)

	startYear, startMonth, startDay := startResult.Date()
	endYear, endMonth, endDay := endResult.Date()
	return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, tz).Unix(), time.Date(endYear, endMonth, endDay, 23, 59, 59, 59, tz).Unix()
}
