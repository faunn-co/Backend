package get_referrals_list

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

type GetReferralList struct {
	c        echo.Context
	ctx      context.Context
	req      *pb.GetReferralListRequest
	key      string
	userId   int64
	userRole int64
}

func New(c echo.Context) *GetReferralList {
	g := new(GetReferralList)
	g.c = c
	g.key = "get_referral_list"
	g.ctx = logger.NewCtx(g.c)
	g.userId = auth_middleware.GetUserIdFromToken(g.c)
	g.userRole = auth_middleware.GetUserRoleFromToken(g.c)
	logger.Info(g.ctx, "GetReferralList Initialized")
	return g
}

func (g *GetReferralList) GetReferralListImpl() ([]*pb.ReferralBasic, *int64, *int64, *resp.Error) {
	if err := g.verifyGetReferralList(); err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var l []*pb.ReferralBasic

	//Filtered for affiliate, return all for admin
	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	k := fmt.Sprintf("%v:%v:%v:%v:%v", g.key, g.userId, g.req.GetTimeSelector().GetPeriod(), start, end)

	if g.userRole == int64(pb.UserRole_ROLE_AFFILIATE) {
		if r := g.cacheGet(k, end); r != nil {
			return r, proto.Int64(start), proto.Int64(end), nil
		}
		if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateReferralListQuery(), start, end, g.userId).Scan(&l).Error; err != nil {
			return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		g.cacheSet(k, l, end)
	} else {
		if g.req.AffiliateName != nil && g.req.GetAffiliateName() != "" {
			//No cache
			if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateReferralListWithNameQuery(), start, end, fmt.Sprintf("%%%v%%", g.req.GetAffiliateName())).Scan(&l).Error; err != nil {
				return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
			}
		} else {
			if r := g.cacheGet(k, end); r != nil {
				return r, proto.Int64(start), proto.Int64(end), nil
			}
			if err := orm.DbInstance(g.ctx).Raw(orm.GetAllReferralListQuery(), start, end).Scan(&l).Error; err != nil {
				return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
			}
			g.cacheSet(k, l, end)
		}
	}
	return l, proto.Int64(start), proto.Int64(end), nil
}

func (g *GetReferralList) verifyGetReferralList() error {
	g.req = new(pb.GetReferralListRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}
	if err := utils.VerifyTimeSelectorFields(g.req.TimeSelector); err != nil {
		return err
	}
	return nil
}

func (g *GetReferralList) cacheGet(k string, end int64) []*pb.ReferralBasic {
	//Read from cache if it is past data, not range, and not today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if val, err := orm.GET(g.c, g.ctx, k, true); err != nil {
			logger.Error(g.ctx, err)
			return nil
		} else if val != nil {
			var redisResp []*pb.ReferralBasic
			jsonErr := json.Unmarshal(val, &redisResp)
			if jsonErr != nil {
				logger.Warn(g.ctx, "GetReferralList | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
			} else {
				logger.Info(g.ctx, "GetReferralList | Successful | Cached %v", k)
				return redisResp
			}
		}
	}
	logger.Info(g.ctx, "GetReferralList | Cache Miss | %v", k)
	return nil
}

func (g *GetReferralList) cacheSet(k string, l []*pb.ReferralBasic, end int64) {
	//Set to cache if it is not period, not within today
	if g.req.GetTimeSelector().GetPeriod() != int64(pb.TimeSelectorPeriod_PERIOD_RANGE) && !utils.IsToday(end) {
		if err := orm.SET(g.ctx, k, l, time.Hour); err != nil {
			logger.ErrorMsg(g.ctx, "GetReferralList | Error while writing to redis: %v", err.Error())
		}
		logger.Info(g.ctx, "GetReferralList | Successful | Written %v to redis", k)
	}
}
