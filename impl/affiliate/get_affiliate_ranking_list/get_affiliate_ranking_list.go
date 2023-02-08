package get_affiliate_ranking_list

import (
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetAffiliateRankingList struct {
	c echo.Context
}

func New(c echo.Context) *GetAffiliateRankingList {
	g := new(GetAffiliateRankingList)
	g.c = c
	return g
}

func (g *GetAffiliateRankingList) GetAffiliateRankingListImpl() (*pb.AffiliateRanking, *pb.AffiliateRanking, *resp.Error) {
	var (
		//Current Period
		r []*pb.AffiliateMetaTopReferral
		c []*pb.AffiliateMetaTopCommission

		//Previous Period
		rP []*pb.AffiliateMetaTopReferral
		cP []*pb.AffiliateMetaTopCommission
	)

	start, end := utils.WeekStartEndDate(time.Now().Unix())
	end = utils.Min(end, time.Now().Unix())
	prevStart, prevEnd := start-utils.WEEK, end-utils.WEEK

	if err := orm.DbInstance(g.c).Raw(orm.Sql6(), start, end).Scan(&r).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	if err := orm.DbInstance(g.c).Raw(orm.Sql7(), start, end).Scan(&c).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	for _, curr := range r {
		var prev *pb.AffiliateMetaTopReferral
		if err := orm.DbInstance(g.c).Raw(orm.Sql8(), curr.GetUserId(), prevStart, prevEnd).Scan(&prev).Error; err != nil {
			return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		if prev == nil {
			curr.PreviousCycleReferrals = proto.Int64(0)
		} else {
			curr.PreviousCycleReferrals = prev.PreviousCycleReferrals
		}
	}

	for _, curr := range c {
		var prev *pb.AffiliateMetaTopCommission
		if err := orm.DbInstance(g.c).Raw(orm.Sql9(), curr.GetUserId(), prevStart, prevEnd).Scan(&prev).Error; err != nil {
			return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
		if prev == nil {
			curr.PreviousCycleCommission = proto.Int64(0)
		} else {
			curr.PreviousCycleCommission = prev.PreviousCycleCommission
		}
	}

	if err := orm.DbInstance(g.c).Raw(orm.Sql6(), prevStart, prevEnd).Scan(&rP).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if err := orm.DbInstance(g.c).Raw(orm.Sql7(), prevStart, prevEnd).Scan(&cP).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return &pb.AffiliateRanking{
			TopAffiliateReferralList:   r,
			TopAffiliateCommissionList: c,
			StartTime:                  proto.Int64(start),
			EndTime:                    proto.Int64(end),
		}, &pb.AffiliateRanking{
			TopAffiliateReferralList:   rP,
			TopAffiliateCommissionList: cP,
			StartTime:                  proto.Int64(prevStart),
			EndTime:                    proto.Int64(prevEnd),
		}, nil
}
