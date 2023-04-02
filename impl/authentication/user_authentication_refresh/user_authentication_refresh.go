package user_authentication_refresh

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type UserAuthenticationRefresh struct {
	c   echo.Context
	ctx context.Context
	req *pb.UserAuthenticationRefreshRequest
}

func New(c echo.Context) *UserAuthenticationRefresh {
	u := new(UserAuthenticationRefresh)
	u.c = c
	u.ctx = logger.NewCtx(u.c)
	logger.Info(u.ctx, "UserAuthenticationRefresh Initialized")
	return u
}

func (u *UserAuthenticationRefresh) UserAuthenticationRefreshImpl() (*pb.Tokens, *resp.Error) {
	if err := u.verifyUserAuthenticationRefresh(); err != nil {
		logger.Error(u.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	if t, err := auth_middleware.Refresh(u.ctx, u.c); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	} else {
		return &pb.Tokens{
			AccessToken:  proto.String(t.AccessToken),
			RefreshToken: proto.String(t.RefreshToken),
		}, nil
	}
}

func (u *UserAuthenticationRefresh) verifyUserAuthenticationRefresh() error {
	u.req = new(pb.UserAuthenticationRefreshRequest)
	if err := u.c.Bind(u.req); err != nil {
		return err
	}
	if u.req.Tokens == nil {
		return errors.New("token is required")
	}
	if u.req.GetTokens().RefreshToken == nil {
		return errors.New("refresh_token is required")
	}
	return nil
}
