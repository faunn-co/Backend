package get_affiliate_list

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

var (
	testAccounts = []string{"test_affiliate"}
)

type GetAffiliateList struct {
	c   echo.Context
	ctx context.Context
	req *pb.GetAffiliateListRequest
	key string
}

func New(c echo.Context) *GetAffiliateList {
	g := new(GetAffiliateList)
	g.c = c
	g.key = "get_affiliate_list"
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateList Initialized")
	return g
}

func (g *GetAffiliateList) GetAffiliateListImpl() ([]*pb.AffiliateMeta, *int64, *int64, *resp.Error) {
	if err := g.verifyGetAffiliateList(); err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	tokenAuth, err := auth_middleware.ExtractTokenMetadata(g.ctx, g.c.Request())
	if err != nil {
		logger.Error(g.ctx, err)
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_TOKEN_ERROR)
	}
	id := tokenAuth.UserId

	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	//cache key : <function>:<user_id>:<period>:<start_ts>:<end_ts>
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, id, g.req.GetTimeSelector().GetPeriod(), start, end)

	if list := g.cacheGet(k, end); list != nil {
		return list, proto.Int64(start), proto.Int64(end), nil
	}
	var l []*pb.AffiliateMeta
	if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateListQuery(), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&l).Error; err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	g.cacheSet(k, l, end)
	return l, proto.Int64(start), proto.Int64(end), nil
}

func (g *GetAffiliateList) verifyGetAffiliateList() error {
	g.req = new(pb.GetAffiliateListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}

func (g *GetAffiliateList) filterTestAccounts(affiliates []*pb.AffiliateMeta) []*pb.AffiliateMeta {
	for i, a := range affiliates {
		for _, t := range testAccounts {
			if a.GetAffiliateName() == t {
				affiliates = append(affiliates[:i], affiliates[i+1:]...)
			}
		}
	}
	return affiliates
}

func (g *GetAffiliateList) cacheGet(k string, end int64) []*pb.AffiliateMeta {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			return nil
		} else if val != nil {
			var redisResp []*pb.AffiliateMeta
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetAffiliateList | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetAffiliateList | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetAffiliateList | Cache Miss | %v", k)
	return nil
}

func (g *GetAffiliateList) cacheSet(k string, l []*pb.AffiliateMeta, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, l, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetAffiliateList | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetAffiliateList | Successful | Written %v to redis", k)
	}
}
