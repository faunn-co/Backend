package get_referral_trend

import (
	"context"
	"database/sql"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetReferralTrend struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetReferralTrendRequest
}

func New(c echo.Context) *GetReferralTrend {
	g := new(GetReferralTrend)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetReferralTrend Initialized")
	return g
}

func (g *GetReferralTrend) GetReferralTrendImpl() ([]*pb.ReferralCoreTimedStats, *resp.Error) {
	if err := g.verifyGetReferralTrend(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	tokenAuth, err := auth_middleware.ExtractTokenMetadata(g.ctx, g.c.Request())
	if err != nil {
		logger.Error(context.Background(), err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_TOKEN_ERROR)
	}

	id := tokenAuth.UserId

	var s []*pb.ReferralCoreTimedStats
	start, end := utils.GetStartEndTimeStampFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralTrendQuery(), sql.Named("id", id), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	//TODO Too expensive, to optimize
	for _, trend := range s {
		type click struct {
			TotalClicks *int64 `json:"total_clicks,omitempty"`
		}
		var c click
		if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralTrendClicksQuery(), sql.Named("id", id), sql.Named("startTime", trend.DateString), sql.Named("endTime", trend.DateString)).Scan(&c).Error; err != nil {
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
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}
