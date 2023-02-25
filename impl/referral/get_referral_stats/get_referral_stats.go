package get_referral_stats

import (
	"context"
	"database/sql"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetReferralStats struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetReferralStatsRequest
}

func New(c echo.Context) *GetReferralStats {
	g := new(GetReferralStats)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetReferralStats Initialized")
	return g
}

func (g *GetReferralStats) GetReferralStatsImpl() (*pb.ReferralStats, *pb.ReferralStats, *resp.Error) {
	if err := g.verifyGetAffiliateStats(); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		//Current Period
		s *pb.ReferralCoreStats
		//Previous Period
		sP *pb.ReferralCoreStats
	)

	start, end, prevStart, prevEnd := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralStatsQuery(), sql.Named("id", g.req.GetAffiliateId()), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&s).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralStatsQuery(), sql.Named("id", g.req.GetAffiliateId()), sql.Named("startTime", prevStart), sql.Named("endTime", prevEnd)).Scan(&sP).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return &pb.ReferralStats{
			CoreStats: s,
			StartTime: proto.Int64(start),
			EndTime:   proto.Int64(end),
		}, &pb.ReferralStats{
			CoreStats: sP,
			StartTime: proto.Int64(prevStart),
			EndTime:   proto.Int64(prevEnd),
		}, nil
}

func (g *GetReferralStats) verifyGetAffiliateStats() error {
	g.req = new(pb.GetReferralStatsRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}
