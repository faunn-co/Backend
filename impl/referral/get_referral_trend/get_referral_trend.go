package get_referral_trend

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetReferralTrend struct {
	c   echo.Context
	req *pb.GetReferralTrendRequest
}

func New(c echo.Context) *GetReferralTrend {
	g := new(GetReferralTrend)
	g.c = c
	return g
}

func (g *GetReferralTrend) GetReferralTrendImpl() ([]*pb.ReferralCoreTimedStats, *resp.Error) {
	if err := g.verifyGetReferralTrend(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		s     []*pb.ReferralCoreTimedStats
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
		startTs, _ = utils.DayStartEndDate(startTs)
		end = utils.ConvertTimeStampYearMonthDay(endTs)
		start = utils.ConvertTimeStampYearMonthDay(startTs)
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS):
		endTs := g.req.GetTimeSelector().GetBaseTs()
		startTs := endTs - utils.MONTH
		startTs, _ = utils.DayStartEndDate(startTs)
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

	if err := orm.DbInstance(g.c).Raw(orm.GetReferralTrendQuery(), sql.Named("id", g.req.GetAffiliateId()), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	fmt.Println(s)

	//TODO Too expensive, to optimize
	for _, trend := range s {
		type click struct {
			TotalClicks *int64 `json:"total_clicks,omitempty"`
		}
		var c click
		if err := orm.DbInstance(g.c).Raw(orm.GetReferralTrendClicksQuery(), sql.Named("id", g.req.GetAffiliateId()), sql.Named("startTime", trend.DateString), sql.Named("endTime", trend.DateString)).Scan(&c).Error; err != nil {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		trend.TotalClicks = c.TotalClicks
		trend.DateString = proto.String(utils.TrimDateString(trend.GetDateString()))
	}
	return s, nil
}

func (g *GetReferralTrend) verifyGetReferralTrend() error {
	g.req = new(pb.GetReferralTrendRequest)
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
