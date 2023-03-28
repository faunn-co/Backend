package user_verification

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

type UserVerification struct {
	c            echo.Context
	ctx          context.Context
	userIdKey    string
	userNameKey  string
	userEmailKey string
}

func New(c echo.Context, ctx context.Context) *UserVerification {
	u := new(UserVerification)
	u.c = c
	u.ctx = ctx
	u.userIdKey = "user_id"
	u.userNameKey = "user_name"
	u.userEmailKey = "user_email"
	return u
}

// VerifyUserId verifies if a certain user_id exists.
// Results are cached indefinitely.
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

	k := fmt.Sprintf("%v:%v", u.userIdKey, userId)
	if val, err := orm.GET(u.c, u.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(u.ctx, "VerifyUserId | Successful | Cached %v", k)
		return nil
	}

	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserIdQuery(), userId).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if err := orm.SET(u.ctx, k, user, time.Hour); err != nil {
		logger.Error(u.ctx, err)
		return nil
	}
	return nil
}

// VerifyUserName verifies if a certain user_name exists.
// Results are cached indefinitely.
func (u *UserVerification) VerifyUserName(name string) error {
	if name == "" {
		return nil
	}

	k := fmt.Sprintf("%v:%v", u.userNameKey, name)
	if val, err := orm.GET(u.c, u.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(u.ctx, "VerifyUserName | Successful | Cached %v", k)
		return nil
	}

	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserNameQuery(), name).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if err := orm.SET(u.ctx, k, user, time.Hour); err != nil {
		logger.Error(u.ctx, err)
		return nil
	}
	return nil
}

// VerifyUserEmail verifies if a certain user_email exists.
// Results are cached indefinitely.
func (u *UserVerification) VerifyUserEmail(email string) error {
	if email == "" {
		return nil
	}

	k := fmt.Sprintf("%v:%v", u.userEmailKey, email)
	if val, err := orm.GET(u.c, u.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(u.ctx, "VerifyUserEmail | Successful | Cached %v", k)
		return nil
	}
	var user *pb.User
	if err := orm.DbInstance(u.ctx).Raw(orm.GetUserInfoWithUserEmailQuery(), email).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if err := orm.SET(u.ctx, k, user, time.Hour); err != nil {
		logger.Error(u.ctx, err)
		return nil
	}
	return nil
}
