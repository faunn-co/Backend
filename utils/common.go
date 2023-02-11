package utils

import (
	"errors"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"strconv"
	"time"
)

func Min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func VerifyTimeSelectorFields(t *pb.TimeSelector) error {
	if t == nil {
		return errors.New("time_selector is required")
	}

	if t.Period == nil {
		return errors.New("period is required")
	}

	if t.GetPeriod() == int64(pb.TimeSelectorPeriod_PERIOD_NONE) {
		return errors.New("invalid period")
	}

	if t.BaseTs == nil && (t.StartTs == nil && t.EndTs == nil) {
		return errors.New("at least base_ts or either of start_ts / end_ts is required")
	}

	if t.StartTs == nil && t.EndTs != nil {
		return errors.New("start_ts is required")
	}

	if t.EndTs == nil && t.StartTs != nil {
		return errors.New("end_ts is required")
	}
	return nil
}

func GetStartEndTimeFromTimeSelector(t *pb.TimeSelector) (int64, int64, int64, int64) {
	var (
		start     int64
		end       int64
		prevStart int64
		prevEnd   int64
	)
	switch t.GetPeriod() {
	case int64(pb.TimeSelectorPeriod_PERIOD_DAY):
		start, end = DayStartEndDate(t.GetBaseTs())
		prevStart, prevEnd = start-DAY, start-SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_WEEK):
		start, end = WeekStartEndDate(t.GetBaseTs())
		prevStart, prevEnd = start-WEEK, start-SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_MONTH):
		//start end of month of ts
		start, end = MonthStartEndDate(t.GetBaseTs())
		prevStart, prevEnd = start-MONTH, start-SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS):
		end = t.GetBaseTs()
		start = end - WEEK
		start, _ = DayStartEndDate(start)
		prevStart, prevEnd = start-WEEK, start-SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS):
		end = t.GetBaseTs()
		start = end - MONTH
		start, _ = DayStartEndDate(start)
		prevStart, prevEnd = start-MONTH, start-SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_RANGE):
		start, _ = DayStartEndDate(t.GetStartTs())
		_, end = DayStartEndDate(t.GetEndTs())
		break
	}
	end = Min(end, time.Now().Unix())
	return start, end, prevStart, prevEnd
}

func GetStartEndTimeStampFromTimeSelector(t *pb.TimeSelector) (string, string) {
	var (
		start string
		end   string
	)
	switch t.GetPeriod() {
	case int64(pb.TimeSelectorPeriod_PERIOD_DAY):
		//start end of ts
		startTs, endTs := DayStartEndDate(t.GetBaseTs())
		endTs = Min(endTs, time.Now().Unix())

		start = ConvertTimeStampYearMonthDay(startTs)
		end = ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_WEEK):
		//start end of week of ts
		startTs, endTs := WeekStartEndDate(t.GetBaseTs())
		endTs = Min(endTs, time.Now().Unix())

		start = ConvertTimeStampYearMonthDay(startTs)
		end = ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_MONTH):
		//start end of month of ts
		startTs, endTs := MonthStartEndDate(t.GetBaseTs())
		endTs = Min(endTs, time.Now().Unix())

		start = ConvertTimeStampYearMonthDay(startTs)
		end = ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS):
		endTs := t.GetBaseTs()
		startTs := endTs - WEEK
		startTs, _ = DayStartEndDate(startTs)
		end = ConvertTimeStampYearMonthDay(endTs)
		start = ConvertTimeStampYearMonthDay(startTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS):
		endTs := t.GetBaseTs()
		startTs := endTs - MONTH
		startTs, _ = DayStartEndDate(startTs)
		end = ConvertTimeStampYearMonthDay(endTs)
		start = ConvertTimeStampYearMonthDay(startTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_RANGE):
		//start of start ts
		//end of end ts
		start = ConvertTimeStampYearMonthDay(t.GetStartTs())
		end = ConvertTimeStampYearMonthDay(t.GetEndTs())
		break
	}
	return start, end
}

func GetStartEndTimeFromPeriod(p string) (int64, int64, int64, int64) {
	var (
		start     int64
		end       int64
		prevStart int64
		prevEnd   int64
	)
	switch p {
	case strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_WEEK), 10):
		start, end = WeekStartEndDate(time.Now().Unix())
		prevStart, prevEnd = start-WEEK, start-SECOND
		break
	case strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_MONTH), 10):
		fallthrough
	default:
		start, end = MonthStartEndDate(time.Now().Unix())
		prevStart, prevEnd = start-MONTH, start-SECOND
		break
	}
	end = Min(end, time.Now().Unix())
	return start, end, prevStart, prevEnd
}
