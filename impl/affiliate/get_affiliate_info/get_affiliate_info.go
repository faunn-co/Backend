package get_affiliate_info

import (
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetAffiliateInfo struct {
	c echo.Context
}

func New(c echo.Context) *GetAffiliateInfo {
	g := new(GetAffiliateInfo)
	g.c = c
	return g
}

func (g *GetAffiliateInfo) GetAffiliateInfoImpl() (*pb.AffiliateMeta, *resp.Error) {
	id := g.c.Param("id")
	if _, err := user_verification.New(g.c).VerifyUserId(id); err != nil {
		return nil, err
	}

	var a *pb.AffiliateMeta
	if err := orm.DbInstance(g.c).Raw(orm.GetAffiliateInfoQuery(), id).Scan(&a).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return a, nil
}
