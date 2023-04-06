package get_user_info

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"time"
)

type GetUserInfo struct {
	c        echo.Context
	ctx      context.Context
	userId   int64
	userRole int64
}

func New(c echo.Context) *GetUserInfo {
	g := new(GetUserInfo)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	g.userId = auth_middleware.GetUserIdFromToken(g.c)
	g.userRole = auth_middleware.GetUserRoleFromToken(g.c)
	logger.Info(g.ctx, "GetUserInfo Initialized")
	return g
}

func (g *GetUserInfo) GetUserInfoImpl() (*pb.AffiliateProfileMeta, *pb.User, *resp.Error) {
	if err := user_verification.New(g.c, g.ctx).VerifyUserId(g.userId); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_NOT_FOUND)
	}

	var a *pb.AffiliateProfileMeta
	t := time.Now().Unix()
	if g.userRole == int64(pb.UserRole_ROLE_ADMIN) {
		if err := orm.DbInstance(g.ctx).Raw(orm.GetAdminInfoQuery(), t, t, t).Scan(&a).Error; err != nil {
			return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
	} else {
		if err := orm.DbInstance(g.ctx).Raw(orm.GetAffiliateInfoQuery(), g.userId, t, t, t).Scan(&a).Error; err != nil {
			return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
	}

	var u *pb.User
	if err := orm.DbInstance(g.ctx).Raw(orm.GetUserInfoWithUserIdQuery(), g.userId).Scan(&u).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	return a, u, nil
}
