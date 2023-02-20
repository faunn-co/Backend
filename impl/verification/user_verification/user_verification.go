package user_verification

import (
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type UserVerification struct {
	c echo.Context
}

func New(c echo.Context) *UserVerification {
	u := new(UserVerification)
	u.c = c
	return u
}

func (u *UserVerification) UserVerificationImpl(id int64) (*pb.User, *resp.Error) {
	var user *pb.User
	if err := orm.DbInstance(u.c).Raw().Scan(&user).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return user, nil
}
