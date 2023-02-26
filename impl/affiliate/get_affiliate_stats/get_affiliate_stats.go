package get_affiliate_stats

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetAffiliateStats struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetAffiliateStatsRequest
}

func New(c echo.Context) *GetAffiliateStats {
	g := new(GetAffiliateStats)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateStats Initialized")
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
	)

	start, end, prevStart, prevEnd := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql5(), start, end).Scan(&s).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql5(), prevStart, prevEnd).Scan(&sP).Error; err != nil {
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
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}
