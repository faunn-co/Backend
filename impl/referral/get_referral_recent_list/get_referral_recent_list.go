package get_referral_recent_list

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

type GetReferralRecentList struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetReferralRecentListRequest
}

func New(c echo.Context) *GetReferralRecentList {
	g := new(GetReferralRecentList)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetReferralRecentList Initialized")
	return g
}

func (g *GetReferralRecentList) GetReferralRecentListImpl() (*pb.ReferralRecent, *resp.Error) {
	if err := g.verifyGetReferralRecentList(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		c []*pb.ReferralClicks
		e []*pb.ReferralEarnings
	)

	start, end, _, _ := utils.GetStartEndTimeFromPeriod(g.c.QueryParam("period"))
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralRecentClicksQuery(), g.req.GetAffiliateId(), start, end).Scan(&c).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralRecentEarningsQuery(), g.req.GetAffiliateId(), start, end).Scan(&e).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return &pb.ReferralRecent{
		RecentClicks:   c,
		RecentEarnings: e,
		StartTime:      proto.Int64(start),
		EndTime:        proto.Int64(end),
	}, nil
}

func (g *GetReferralRecentList) verifyGetReferralRecentList() error {
	g.req = new(pb.GetReferralRecentListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	return nil
}
