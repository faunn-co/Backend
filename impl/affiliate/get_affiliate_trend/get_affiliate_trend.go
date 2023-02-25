package get_affiliate_trend

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
)

type GetAffiliateTrend struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetAffiliateTrendRequest
}

func New(c echo.Context) *GetAffiliateTrend {
	g := new(GetAffiliateTrend)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateTrend Initialized")
	return g
}

func (g *GetAffiliateTrend) GetAffiliateTrendImpl() ([]*pb.AffiliateCoreTimedStats, *resp.Error) {
	if err := g.verifyGetAffiliateTrend(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var s []*pb.AffiliateCoreTimedStats
	start, end := utils.GetStartEndTimeStampFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql10(), start, end).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return s, nil
}

func (g *GetAffiliateTrend) verifyGetAffiliateTrend() error {
	g.req = new(pb.GetAffiliateTrendRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}
