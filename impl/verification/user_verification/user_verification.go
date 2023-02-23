package user_verification

import (
	"errors"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"strconv"
)

type UserVerification struct {
	c echo.Context
}

func New(c echo.Context) *UserVerification {
	u := new(UserVerification)
	u.c = c
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
	if err := orm.DbInstance(u.c).Raw(orm.GetUserInfoQuery(), userId).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return nil
}
