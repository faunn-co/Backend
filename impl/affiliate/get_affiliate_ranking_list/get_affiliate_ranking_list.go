package get_affiliate_ranking_list

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

type GetAffiliateRankingList struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *GetAffiliateRankingList {
	g := new(GetAffiliateRankingList)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateRankingList Initialized")
	return g
}

func (g *GetAffiliateRankingList) GetAffiliateRankingListImpl() (*pb.AffiliateRanking, *resp.Error) {
	var (
		//Current Period
		r         []*pb.AffiliateMetaTopReferral
		c         []*pb.AffiliateMetaTopCommission
		start     int64
		end       int64
		prevStart int64
		prevEnd   int64
		err       error
	)

	p := g.c.QueryParam("period")
	if start, end, prevStart, prevEnd, err = utils.GetStartEndTimeFromPeriod(p); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql6(), start, end).Scan(&r).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.ctx).Raw(orm.Sql7(), start, end).Scan(&c).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	for _, curr := range r {
		var prev *pb.AffiliateMetaTopReferral
		if err := orm.DbInstance(g.ctx).Raw(orm.Sql8(), curr.GetUserId(), prevStart, prevEnd).Scan(&prev).Error; err != nil {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		if prev == nil {
			curr.PreviousCycleReferrals = proto.Int64(0)
		} else {
			curr.PreviousCycleReferrals = prev.PreviousCycleReferrals
		}
	}

	for _, curr := range c {
		var prev *pb.AffiliateMetaTopCommission
		if err := orm.DbInstance(g.ctx).Raw(orm.Sql9(), curr.GetUserId(), prevStart, prevEnd).Scan(&prev).Error; err != nil {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		if prev == nil {
			curr.PreviousCycleCommission = proto.Int64(0)
		} else {
			curr.PreviousCycleCommission = prev.PreviousCycleCommission
		}
	}

	return &pb.AffiliateRanking{
		TopAffiliateReferralList:   r,
		TopAffiliateCommissionList: c,
		StartTime:                  proto.Int64(start),
		EndTime:                    proto.Int64(end),
	}, nil
}
