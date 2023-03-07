package user_deauthentication

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type UserDeAuthentication struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *UserDeAuthentication {
	u := new(UserDeAuthentication)
	u.c = c
	u.ctx = logger.NewCtx(u.c)
	logger.Info(u.ctx, "UserDeAuthentication Initialized")
	return u
}

func (u *UserDeAuthentication) UserDeAuthenticationImpl() *resp.Error {
	if err := u.executeLogout(); err != nil {
		logger.Error(u.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_LOGOUT_FAIL)
	}
	return nil
}

func (u *UserDeAuthentication) executeLogout() error {
	tokenAuth, err := auth_middleware.ExtractTokenMetadata(u.ctx, u.c.Request())
	if err != nil {
		logger.Error(context.Background(), err)
		return err
	}

	deleted, err := auth_middleware.DeleteAuth(u.ctx, tokenAuth.AccessUuid)
	if err != nil {
		return err
	}

	if deleted == 0 {
		logger.Warn(u.ctx, "access_token does not exist")
	}

	return nil
}
