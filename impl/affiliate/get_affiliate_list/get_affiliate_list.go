package get_affiliate_list

import (
	"database/sql"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GetAffiliateList struct {
	c   echo.Context
	req *pb.GetAffiliateListRequest
}

func New(c echo.Context) *GetAffiliateList {
	g := new(GetAffiliateList)
	g.c = c
	return g
}

func (g *GetAffiliateList) GetAffiliateListImpl() ([]*pb.AffiliateMeta, *int64, *int64, *resp.Error) {
	if err := g.verifyGetAffiliateList(); err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var l []*pb.AffiliateMeta
	start, end, _, _ := utils.GetStartEndTimeFromTimeSelector(g.req.GetTimeSelector())
	if err := orm.DbInstance(g.c).Raw(orm.GetAffiliateListQuery(), sql.Named("startTime", start), sql.Named("endTime", end)).Scan(&l).Error; err != nil {
		return nil, nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
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
