package get_affiliate_stats

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
	"google.golang.org/protobuf/proto"
	"time"
)

type GetAffiliateStats struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetAffiliateStatsRequest
	key string
}

func New(c echo.Context) *GetAffiliateStats {
	g := new(GetAffiliateStats)
	g.c = c
	g.key = "get_affiliate_stats"
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateStats Initialized")
	return g
}

func (g *GetAffiliateStats) GetAffiliateStatsImpl() (*pb.GetAffiliateStatsResponse, *resp.Error) {
	if err := g.verifyGetAffiliateStats(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	tokenAuth, err := auth_middleware.ExtractTokenMetadata(g.ctx, g.c.Request())
	if err != nil {
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_TOKEN_ERROR)
	}
	id := tokenAuth.UserId

	var (
		//Current Period
		s *pb.AffiliateCoreStats
		//Previous Period
		sP       *pb.AffiliateCoreStats
		response *pb.GetAffiliateStatsResponse
	)

	start, end, prevStart, prevEnd := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, id, g.req.GetTimeSelector().GetPeriod(), start, end)

	if r := g.cacheGet(k, end); err != nil {
		return r, nil
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql5(), start, end).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql5(), prevStart, prevEnd).Scan(&sP).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	response = &pb.GetAffiliateStatsResponse{
		AffiliateStats: &pb.AffiliateStats{
			CoreStats: s,
			StartTime: proto.Int64(start),
			EndTime:   proto.Int64(end),
		},
		AffiliateStatsPreviousCycle: &pb.AffiliateStats{
			CoreStats: sP,
			StartTime: proto.Int64(prevStart),
			EndTime:   proto.Int64(prevEnd),
		},
	}

	g.cacheSet(k, response, end)
	return response, nil
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

func (g *GetAffiliateStats) cacheGet(k string, end int64) *pb.GetAffiliateStatsResponse {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			return nil
		} else if val != nil {
			var redisResp *pb.GetAffiliateStatsResponse
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetAffiliateStats | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetAffiliateStats | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetAffiliateStats | Cache Miss | %v", k)
	return nil
}

func (g *GetAffiliateStats) cacheSet(k string, response *pb.GetAffiliateStatsResponse, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, response, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetAffiliateStats | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetAffiliateStats | Successful | Written %v to redis", k)
	}
}
