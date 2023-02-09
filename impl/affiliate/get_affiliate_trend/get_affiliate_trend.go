package get_affiliate_trend

import (
	"errors"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"time"
)

type GetAffiliateTrend struct {
	c   echo.Context
	req *pb.GetAffiliateTrendRequest
}

func New(c echo.Context) *GetAffiliateTrend {
	g := new(GetAffiliateTrend)
	g.c = c
	return g
}

func (g *GetAffiliateTrend) GetAffiliateTrendImpl() ([]*pb.AffiliateCoreTimedStats, *resp.Error) {
	if err := g.verifyGetAffiliateTrend(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		s     []*pb.AffiliateCoreTimedStats
		start string
		end   string
	)

	switch g.req.GetTimeSelector().GetPeriod() {
	case int64(pb.TimeSelectorPeriod_PERIOD_DAY):
		//start end of ts
		startTs, endTs := utils.DayStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		endTs = utils.Min(endTs, time.Now().Unix())

		start = utils.ConvertTimeStampYearMonthDay(startTs)
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_WEEK):
		//start end of week of ts
		startTs, endTs := utils.WeekStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		endTs = utils.Min(endTs, time.Now().Unix())

		start = utils.ConvertTimeStampYearMonthDay(startTs)
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_MONTH):
		//start end of month of ts
		startTs, endTs := utils.MonthStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		endTs = utils.Min(endTs, time.Now().Unix())

		start = utils.ConvertTimeStampYearMonthDay(startTs)
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS):
		endTs := g.req.GetTimeSelector().GetBaseTs()
		startTs := endTs - utils.WEEK
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		start = utils.ConvertTimeStampYearMonthDay(startTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS):
		endTs := g.req.GetTimeSelector().GetBaseTs()
		startTs := endTs - utils.MONTH
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		start = utils.ConvertTimeStampYearMonthDay(startTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_RANGE):
		//start of start ts
		//end of end ts
		start = utils.ConvertTimeStampYearMonthDay(g.req.GetTimeSelector().GetStartTs())
		end = utils.ConvertTimeStampYearMonthDay(g.req.GetTimeSelector().GetEndTs())
		break
	}

	if err := orm.DbInstance(g.c).Raw(orm.Sql10(), start, end).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return s, nil
}

func (g *GetAffiliateTrend) verifyGetAffiliateTrend() error {
	g.req = new(pb.GetAffiliateTrendRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if g.req.TimeSelector == nil {
		return errors.New("time_selector is required")
	}

	if g.req.GetTimeSelector().Period == nil {
		return errors.New("period is required")
	}

	if g.req.GetTimeSelector().GetPeriod() == int64(pb.TimeSelectorPeriod_PERIOD_NONE) {
		return errors.New("invalid period")
	}

	if g.req.GetTimeSelector().BaseTs == nil && (g.req.GetTimeSelector().StartTs == nil && g.req.GetTimeSelector().EndTs == nil) {
		return errors.New("at least base_ts or either of start_ts / end_ts is required")
	}

	if g.req.GetTimeSelector().StartTs == nil && g.req.GetTimeSelector().EndTs != nil {
		return errors.New("start_ts is required")
	}

	if g.req.GetTimeSelector().EndTs == nil && g.req.GetTimeSelector().StartTs != nil {
		return errors.New("end_ts is required")
	}
	return nil
}
