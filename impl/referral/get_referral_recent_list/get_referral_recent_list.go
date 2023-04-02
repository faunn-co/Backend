package get_referral_recent_list

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
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

	tokenAuth, err := auth_middleware.ExtractTokenMetadata(g.ctx, g.c.Request())
	if err != nil {
		logger.Error(context.Background(), err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_TOKEN_ERROR)
	}

	id := tokenAuth.UserId

	var (
		c []*pb.ReferralClicks
		e []*pb.ReferralEarnings
	)

	t := time.Now().Unix()
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralRecentClicksQuery(), id, t).Scan(&c).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralRecentEarningsQuery(), id, t).Scan(&e).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return &pb.ReferralRecent{
		RecentClicks:   c,
		RecentEarnings: e,
		EndTime:        proto.Int64(t),
	}, nil
}

func (g *GetReferralRecentList) verifyGetReferralRecentList() error {
	g.req = new(pb.GetReferralRecentListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	return nil
}
