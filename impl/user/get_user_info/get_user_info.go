package get_user_info

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetUserInfo struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *GetUserInfo {
	g := new(GetUserInfo)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetUserInfo Initialized")
	return g
}

func (g *GetUserInfo) GetUserInfoImpl() (*pb.AffiliateProfileMeta, *pb.User, *resp.Error) {
	id := g.c.QueryParam("id")
	fmt.Println(id)
	if err := user_verification.New(g.c, g.ctx).VerifyUserId(id); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_NOT_FOUND)
	}

	var a *pb.AffiliateProfileMeta
	if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateInfoQuery(), id).Scan(&a).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	var u *pb.User
	if err := orm.DbInstance(g.ctx).Raw(orm.GetUserInfoWithUserIdQuery(), id).Scan(&u).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return a, u, nil
}
