package get_affiliate_stats

import (
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetAffiliateStats struct {
	c echo.Context
}

func New(c echo.Context) *GetAffiliateStats {
	g := new(GetAffiliateStats)
	g.c = c
	return g
}

func (g *GetAffiliateStats) GetAffiliateStatsImpl() (*pb.AffiliateStats, *pb.AffiliateStats, *resp.Error) {
	var (
		//Current Period
		s *pb.AffiliateCoreStats

		//Previous Period
		sP *pb.AffiliateCoreStats
	)

	start, end := utils.WeekStartEndDate(time.Now().Unix())
	end = utils.Min(end, time.Now().Unix())
	prevStart, prevEnd := start-utils.WEEK, start-utils.SECOND

	//TODO add time stats filter
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
