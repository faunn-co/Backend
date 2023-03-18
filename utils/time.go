package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	//TimeZone constant timezone
	TimeZone = "Asia/Singapore"
	SECOND   = 1
	MINUTE   = 60 * SECOND
	HOUR     = 60 * MINUTE
	DAY      = 24 * HOUR
	WEEK     = 7 * DAY
	MONTH    = 4 * WEEK
)

var (
	tz, _ = time.LoadLocation(TimeZone)
)

// UnixToUTC Converts current unix time to UTC time object
func UnixToUTC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Local().UTC()
}

func DayStartEndDate(timestamp int64) (int64, int64) {
	year, month, day := UnixToUTC(timestamp).In(tz).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, tz).Unix(), time.Date(year, month, day, 23, 59, 59, 59, tz).Unix()
}

func IsToday(endTime int64) bool {
	start, _ := DayStartEndDate(time.Now().Unix())
	return endTime >= start
}

// WeekStartEndDate Returns the start and end day of the current week in SGT unix time
func WeekStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)

	startOffset := (int(time.Sunday) - int(date.Weekday()) - 7) % 7
	startResult := date.Add(time.Duration(startOffset*24) * time.Hour)
	endResult := startResult.Add(time.Duration(6*24) * time.Hour)

	startYear, startMonth, startDay := startResult.Date()
	endYear, endMonth, endDay := endResult.Date()
	return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, tz).Unix(), time.Date(endYear, endMonth, endDay, 23, 59, 59, 59, tz).Unix()
}

// MonthStartEndDate Returns the start and end day of the current month in SGT unix time
func MonthStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)
	currentYear, currentMonth, _ := date.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, tz)
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 23, 59, 59, 59, tz)
	return firstOfMonth.Unix(), lastOfMonth.Unix()
}

// ConvertTimeStampYearMonthDay Converts current unix timestamp to yyyy-mm-dd format
// time format: Mon Jan 2 15:04:05 -0700 MST 2006
func ConvertTimeStampYearMonthDay(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("2006-01-02"))
}

func TrimDateString(d string) string {
	split := strings.Split(d, "-")
	return fmt.Sprintf("%v-%v", split[1], split[2])
}

func GetDayTSWithDate(date string) (int64, error) {
	splitDate := strings.Split(date, "-")
	if len(splitDate) < 3 {
		return 0, errors.New("invalid date format")
	}
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return 0, err
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return 0, err
	}
	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return 0, err
	}
	m := time.Month(month)
	return time.Date(year, m, day, 0, 0, 0, 0, tz).Unix(), nil
}
