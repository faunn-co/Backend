package get_referral_stats

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

type GetReferralStats struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetReferralStatsRequest
	key string
}

func New(c echo.Context) *GetReferralStats {
	g := new(GetReferralStats)
	g.c = c
	g.key = "get_referral_stats"
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetReferralStats Initialized")
	return g
}

func (g *GetReferralStats) GetReferralStatsImpl() (*pb.GetReferralStatsResponse, *resp.Error) {
	if err := g.verifyGetAffiliateStats(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		//Current Period
		s *pb.ReferralCoreStats
		//Previous Period
		sP *pb.ReferralCoreStats
	)

	tokenAuth, err := auth_middleware.ExtractTokenMetadata(g.ctx, g.c.Request())
	if err != nil {
		logger.Error(context.Background(), err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_TOKEN_ERROR)
	}
	id := tokenAuth.UserId

	start, end, prevStart, prevEnd := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, id, g.req.GetTimeSelector().GetPeriod(), start, end)

	if r := g.cacheGet(k, end); r != nil {
		return r, nil
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralStatsQuery(), sql.Named("id", id), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&s).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralStatsQuery(), sql.Named("id", id), sql.Named("startTime", prevStart), sql.Named("endTime", prevEnd)).Scan(&sP).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	response := &pb.GetReferralStatsResponse{
		ResponseMeta: nil,
		ReferralStats: &pb.ReferralStats{
			CoreStats: s,
			StartTime: proto.Int64(start),
			EndTime:   proto.Int64(end),
		},
		ReferralStatsPreviousCycle: &pb.ReferralStats{
			CoreStats: sP,
			StartTime: proto.Int64(prevStart),
			EndTime:   proto.Int64(prevEnd),
		},
	}
	g.cacheSet(k, response, end)
	return response, nil
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

func (g *GetReferralStats) cacheGet(k string, end int64) *pb.GetReferralStatsResponse {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			return nil
		} else if val != nil {
			var redisResp *pb.GetReferralStatsResponse
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetReferralStats | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetReferralStats | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetReferralStats | Cache Miss | %v", k)
	return nil
}

func (g *GetReferralStats) cacheSet(k string, response *pb.GetReferralStatsResponse, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, response, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetReferralStats | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetReferralStats | Successful | Written %v to redis", k)
	}
}
