package get_affiliate_stats

import (
	"errors"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetAffiliateStats struct {
	c   echo.Context
	req *pb.GetAffiliateStatsRequest
}

func New(c echo.Context) *GetAffiliateStats {
	g := new(GetAffiliateStats)
	g.c = c
	return g
}

func (g *GetAffiliateStats) GetAffiliateStatsImpl() (*pb.AffiliateStats, *pb.AffiliateStats, *resp.Error) {
	if err := g.verifyGetAffiliateStats(); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		//Current Period
		s *pb.AffiliateCoreStats

		//Previous Period
		sP *pb.AffiliateCoreStats

		start     int64
		end       int64
		prevStart int64
		prevEnd   int64
	)

	switch g.req.GetTimeSelector().GetPeriod() {
	case int64(pb.TimeSelectorPeriod_PERIOD_DAY):
		//start end of ts
		start, end = utils.DayStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		prevStart, prevEnd = start-utils.DAY, start-utils.SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_WEEK):
		//start end of week of ts
		start, end = utils.WeekStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		prevStart, prevEnd = start-utils.WEEK, start-utils.SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_MONTH):
		//start end of month of ts
		start, end = utils.MonthStartEndDate(g.req.GetTimeSelector().GetBaseTs())
		prevStart, prevEnd = start-utils.MONTH, start-utils.SECOND
		break
	case int64(pb.TimeSelectorPeriod_PERIOD_RANGE):
		//start of start ts
		//end of end ts
		start, _ = utils.DayStartEndDate(g.req.GetTimeSelector().GetStartTs())
		_, end = utils.DayStartEndDate(g.req.GetTimeSelector().GetEndTs())
		break
	}

	end = utils.Min(end, time.Now().Unix())

	if err := orm.DbInstance(g.c).Raw(orm.Sql5(), start, end).Scan(&s).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.c).Raw(orm.Sql5(), prevStart, prevEnd).Scan(&sP).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return &pb.AffiliateStats{
			CoreStats: s,
			StartTime: proto.Int64(start),
			EndTime:   proto.Int64(end),
		}, &pb.AffiliateStats{
			CoreStats: sP,
			StartTime: proto.Int64(prevStart),
			EndTime:   proto.Int64(prevEnd),
		}, nil
}

func (g *GetAffiliateStats) verifyGetAffiliateStats() error {
	g.req = new(pb.GetAffiliateStatsRequest)
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
