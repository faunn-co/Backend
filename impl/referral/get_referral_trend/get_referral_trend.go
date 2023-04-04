package get_referral_trend

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetReferralTrend struct {
	c      echo.Context
	ctx    context.Context
	req    *pb.GetReferralTrendRequest
	key    string
	userId int64
}

func New(c echo.Context) *GetReferralTrend {
	g := new(GetReferralTrend)
	g.c = c
	g.key = "get_referral_trend"
	g.ctx = logger.NewCtx(g.c)
	g.userId = auth_middleware.GetUserIdFromToken(g.c)
	logger.Info(g.ctx, "GetReferralTrend Initialized")
	return g
}

func (g *GetReferralTrend) GetReferralTrendImpl() ([]*pb.ReferralCoreTimedStats, *resp.Error) {
	if err := g.verifyGetReferralTrend(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var s []*pb.ReferralCoreTimedStats
	_, endTs, start, end := utils.GetStartEndTimeStampFromTimeSelector(g.req.GetTimeSelector())
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, g.userId, g.req.GetTimeSelector().GetPeriod(), start, end)

	if r := g.cacheGet(k, endTs); r != nil {
		return r, nil
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralTrendQuery(), sql.Named("id", g.userId), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	//TODO Too expensive, to optimize
	for _, trend := range s {
		type click struct {
			TotalClicks *int64 `json:"total_clicks,omitempty"`
		}
		var c click
		if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralTrendClicksQuery(), sql.Named("id", g.userId), sql.Named("startTime", trend.DateString), sql.Named("endTime", trend.DateString)).Scan(&c).Error; err != nil {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		trend.TotalClicks = c.TotalClicks
		trend.DateString = proto.String(utils.TrimDateString(trend.GetDateString()))
	}

	g.cacheSet(k, s, endTs)
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

func (g *GetReferralTrend) cacheGet(k string, end int64) []*pb.ReferralCoreTimedStats {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			logger.Error(g.ctx, err)
			return nil
		} else if val != nil {
			var redisResp []*pb.ReferralCoreTimedStats
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetReferralTrend | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetReferralTrend | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetReferralTrend | Cache Miss | %v", k)
	return nil
}

func (g *GetReferralTrend) cacheSet(k string, stats []*pb.ReferralCoreTimedStats, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, stats, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetReferralTrend | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetReferralTrend | Successful | Written %v to redis", k)
	}
}
