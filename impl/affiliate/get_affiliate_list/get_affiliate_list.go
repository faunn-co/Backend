package get_affiliate_list

import (
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetAffiliateList struct {
	c echo.Context
}

func New(c echo.Context) *GetAffiliateList {
	g := new(GetAffiliateList)
	g.c = c
	return g
}

func (g *GetAffiliateList) GetAffiliateListImpl() ([]*pb.AffiliateMeta, *resp.Error) {
	var list []*pb.AffiliateMeta
	if err := orm.DbInstance(g.c).Raw(orm.Sql1()).Scan(&list).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	fmt.Println(list)
	return list, nil
}
