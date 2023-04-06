package get_affiliate_trend

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"time"
)

type GetAffiliateTrend struct {
	c      echo.Context
	ctx    context.Context
	req    *pb.GetAffiliateTrendRequest
	key    string
	userId int64
}

func New(c echo.Context) *GetAffiliateTrend {
	g := new(GetAffiliateTrend)
	g.c = c
	g.key = "get_affiliate_trend"
	g.ctx = logger.NewCtx(g.c)
	g.userId = auth_middleware.GetUserIdFromToken(g.c)
	logger.Info(g.ctx, "GetAffiliateTrend Initialized")
	return g
}

func (g *GetAffiliateTrend) GetAffiliateTrendImpl() ([]*pb.AffiliateCoreTimedStats, *resp.Error) {
	if err := g.verifyGetAffiliateTrend(); err != nil {
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var s []*pb.AffiliateCoreTimedStats
	_, endTs, start, end := utils.GetStartEndTimeStampFromTimeSelector(g.req.GetTimeSelector())
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, g.userId, g.req.GetTimeSelector().GetPeriod(), start, end)

	if r := g.cacheGet(k, endTs); r != nil {
		return r, nil
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql10(), start, end).Scan(&s).Error; err != nil {
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	g.cacheSet(k, s, endTs)
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

func (g *GetAffiliateTrend) cacheGet(k string, end int64) []*pb.AffiliateCoreTimedStats {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			return nil
		} else if val != nil {
			var redisResp []*pb.AffiliateCoreTimedStats
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetAffiliateTrend | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetAffiliateTrend | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetAffiliateTrend | Cache Miss | %v", k)
	return nil
}

func (g *GetAffiliateTrend) cacheSet(k string, stats []*pb.AffiliateCoreTimedStats, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, stats, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetAffiliateTrend | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetAffiliateTrend | Successful | Written %v to redis", k)
	}
}
