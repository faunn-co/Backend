package get_affiliate_list

import (
	"context"
	"database/sql"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
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
	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())

	//TODO Cache time slot specific
	//if val, err := orm.GET(g.c, g.key); err != nil {
	//	return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_REDIS)
	//} else if val != nil {
	//	var redisResp []*pb.AffiliateMeta
	//	jsonErr := json.Unmarshal(val, &redisResp)
	//	if jsonErr != nil {
	//		log.Warnf("GetAffiliateList | Fail to unmarshal Redis value of key %v : %v, reading from API", g.key, jsonErr)
	//	} else {
	//		log.Infof("GetAffiliateList | Successful | Cached %v", g.key)
	//		return redisResp, proto.Int64(start), proto.Int64(end), nil
	//	}
	//}

	var l []*pb.AffiliateMeta
	if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateListQuery(), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&l).Error; err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	//if err := orm.SET(g.c, g.key, l, time.Hour); err != nil {
	//	log.Errorf("GetAffiliateList | Error while writing to redis: %v", err.Error())
	//	return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_REDIS)
	//}
	//log.Infof("GetAffiliateList | Successful | Written %v to redis", g.key)
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
