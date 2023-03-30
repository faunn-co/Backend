package user_authentication

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/encrypt"
	"github.com/aaronangxz/AffiliateManager/impl/verification/user_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type UserAuthentication struct {
	c   echo.Context
	ctx context.Context
	req *pb.UserAuthenticationRequest
}

func New(c echo.Context) *UserAuthentication {
	u := new(UserAuthentication)
	u.c = c
	u.ctx = logger.NewCtx(u.c)
	logger.Info(u.ctx, "UserAuthentication Initialized")
	return u
}

func (u *UserAuthentication) UserAuthenticationImpl() (*pb.AuthCookie, *resp.Error) {
	if err := u.verifyUserAuthentication(); err != nil {
		logger.Error(u.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	if err := u.verifyUserEmail(); err != nil {
		logger.Error(u.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_USER_NOT_FOUND)
	}

	var cookie *pb.AuthCookie
	var err error
	if cookie, err = u.executeLogin(); err != nil {
		logger.Error(u.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_LOGIN_FAIL)
	}
	return cookie, nil
}

func (u *UserAuthentication) executeLogin() (*pb.AuthCookie, error) {
	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithAuthQuery(), u.req.GetUserName()).Scan(&user).Error; err != nil {
		return nil, err
	}

	if !encrypt.ComparePasswords(u.ctx, user.GetUserPassword(), u.req.GetUserPassword()) {
		return nil, errors.New("login credentials are incorrect")
	}

	token, err := auth_middleware.CreateToken(u.ctx, user.GetUserId(), user.GetUserRole())
	if err != nil {
		return nil, err
	}

	saveErr := auth_middleware.CreateAuth(u.ctx, user.GetUserId(), token)
	if saveErr != nil {
		return nil, saveErr
	}

	return &pb.AuthCookie{
		UserId:    user.UserId,
		UserName:  user.UserName,
		UserEmail: user.UserEmail,
		UserRole:  user.UserRole,
		Tokens: &pb.Tokens{
			AccessToken:  proto.String(token.AccessToken),
			RefreshToken: proto.String(token.RefreshToken),
		},
	}, nil
}

func (u *UserAuthentication) verifyUserEmail() error {
	if err := user_verification.New(u.c, u.ctx).VerifyUserName(u.req.GetUserName()); err != nil {
		return err
	}
	return nil
}

func (u *UserAuthentication) verifyUserAuthentication() error {
	u.req = new(pb.UserAuthenticationRequest)
	if err := u.c.Bind(u.req); err != nil {
		return err
	}
	if u.req == nil {
		return errors.New("request body is empty")
	}

	if u.req.UserName == nil {
		return errors.New("username is empty")
	}

	if u.req.UserPassword == nil {
		return errors.New("password is empty")
	}

	return nil
}
