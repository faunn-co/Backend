package get_referrals_list

import (
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetReferralList struct {
	c   echo.Context
	req *pb.GetReferralListRequest
}

func New(c echo.Context) *GetReferralList {
	g := new(GetReferralList)
	g.c = c
	return g
}

func (g *GetReferralList) GetReferralListImpl() ([]*pb.ReferralBasic, *int64, *int64, *resp.Error) {
	if err := g.verifyGetReferralList(); err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	var l []*pb.ReferralBasic
	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())

	if g.req.AffiliateId != nil {
		if err := orm.DbInstance(g.c).Raw(orm.GetAffiliateReferralListQuery(), start, end, g.req.GetAffiliateId()).Scan(&l).Error; err != nil {
			return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
	} else {
		if g.req.AffiliateName != nil {
			if err := orm.DbInstance(g.c).Raw(orm.GetAffiliateReferralListWithNameQuery(), start, end, fmt.Sprintf("%%%v%%", g.req.GetAffiliateName())).Scan(&l).Error; err != nil {
				return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
			}
		} else {
			if err := orm.DbInstance(g.c).Raw(orm.GetAllReferralListQuery(), start, end).Scan(&l).Error; err != nil {
				return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
			}
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
