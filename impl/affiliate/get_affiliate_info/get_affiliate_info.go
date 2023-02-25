package get_affiliate_info

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetAffiliateInfo struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *GetAffiliateInfo {
	g := new(GetAffiliateInfo)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetAffiliateInfo Initialized")
	return g
}

func (g *GetAffiliateInfo) GetAffiliateInfoImpl() (*pb.AffiliateMeta, *resp.Error) {
	id := g.c.Param("id")
	if err := user_verification.New(g.c, g.ctx).VerifyUserId(id); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_NOT_FOUND)
	}

	var a *pb.AffiliateMeta
	if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateInfoQuery(), id).Scan(&a).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return a, nil
}
