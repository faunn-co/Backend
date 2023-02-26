package user_verification

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"strconv"
)

type UserVerification struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context, ctx context.Context) *UserVerification {
	u := new(UserVerification)
	u.c = c
	u.ctx = ctx
	return u
}

func (u *UserVerification) VerifyUserId(id interface{}) error {
	var userId int64
	switch id.(type) {
	case string:
		var err error
		userId, err = strconv.ParseInt(id.(string), 10, 64)
		if err != nil {
			return err
		}
		break
	case int64:
		userId = id.(int64)
		break
	}

	if userId == 0 {
		return nil
	}
	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserIdQuery(), userId).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return nil
}

func (u *UserVerification) VerifyUserName(name string) error {
	if name == "" {
		return nil
	}
	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserNameQuery(), name).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return nil
}

func (u *UserVerification) VerifyUserEmail(email string) error {
	if email == "" {
		return nil
	}
	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserEmailQuery(), email).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return nil
}
