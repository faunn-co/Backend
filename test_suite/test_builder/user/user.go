package user

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/encrypt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/test_utils"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

type User struct {
	UserInfo      *pb.User
	AffiliateInfo *pb.AffiliateDetailsDb
	Token         *pb.Tokens
	token         *auth_middleware.TokenDetails
	NeedLogin     bool
}

func New() *User {
	u := new(User)
	u.UserInfo = new(pb.User)
	u.AffiliateInfo = new(pb.AffiliateDetailsDb)
	u.Token = new(pb.Tokens)
	u.NeedLogin = true
	orm.ENV = "TEST"
	return u
}

func (u *User) SetNeedLogin(needLogin bool) *User {
	u.NeedLogin = needLogin
	return u
}

func (u *User) SetUserName(name string) *User {
	u.UserInfo.UserName = proto.String(name)
	return u
}

func (u *User) SetUserEmail(email string) *User {
	u.UserInfo.UserEmail = proto.String(email)
	return u
}

func (u *User) SetUserContact(contact string) *User {
	u.UserInfo.UserContact = proto.String(contact)
	return u
}

func (u *User) SetUserRole(role int64) *User {
	u.UserInfo.UserRole = proto.Int64(role)
	return u
}

func (u *User) SetEntityName(name string) *User {
	u.AffiliateInfo.EntityName = proto.String(name)
	return u
}

func (u *User) SetEntityIdentifier(identifier string) *User {
	u.AffiliateInfo.EntityIdentifier = proto.String(identifier)
	return u
}

func (u *User) SetAffiliateType(affiliateType int64) *User {
	u.AffiliateInfo.AffiliateType = proto.Int64(affiliateType)
	return u
}

func (u *User) SetUniqueReferralCode(code string) *User {
	u.AffiliateInfo.UniqueReferralCode = proto.String(code)
	return u
}

func (u *User) filDefaults() *User {
	if u.UserInfo.UserName == nil {
		u.SetUserName(test_utils.RandomStringWithCharset(10))
	}

	if u.UserInfo.UserEmail == nil {
		u.SetUserEmail(fmt.Sprintf("%v@random.com", u.UserInfo.GetUserName()))
	}

	if u.UserInfo.UserContact == nil {
		u.SetUserContact(fmt.Sprintf("+60%v", test_utils.RandomRange(1000000, 9999999)))
	}

	if u.UserInfo.UserRole == nil {
		u.SetUserRole(int64(pb.UserRole_ROLE_AFFILIATE))
	}

	if u.AffiliateInfo.EntityName == nil {
		u.SetEntityName(test_utils.RandomStringWithCharset(10))
	}

	if u.AffiliateInfo.EntityIdentifier == nil {
		u.SetEntityIdentifier(test_utils.RandomStringWithCharset(15))
	}

	if u.AffiliateInfo.AffiliateType == nil {
		if u.UserInfo.UserRole != nil && u.UserInfo.GetUserRole() == int64(pb.UserRole_ROLE_AFFILIATE) {
			u.SetAffiliateType(test_utils.RandomRange(0, 1))
		}
	}

	if u.AffiliateInfo.UniqueReferralCode == nil {
		u.SetUniqueReferralCode(test_utils.RandomStringWithCharset(5))
	}
	return u
}

func (u *User) Build() *User {
	u.filDefaults()

	type User struct {
		UserId          int64 `gorm:"primary_key"`
		UserName        string
		UserEmail       string
		UserContact     string
		UserRole        int64
		CreateTimestamp int64
	}

	user := User{
		UserName:        *u.UserInfo.UserName,
		UserEmail:       *u.UserInfo.UserEmail,
		UserContact:     *u.UserInfo.UserContact,
		UserRole:        *u.UserInfo.UserRole,
		CreateTimestamp: time.Now().Unix(),
	}

	if err := orm.DbInstance(context.Background()).Table(orm.USER_TABLE).Create(&user).Error; err != nil {
		log.Error(err)
		return nil
	}

	type UserAuth struct {
		UserId       int64 `gorm:"primary_key"`
		UserPassword string
	}

	auth := UserAuth{
		UserId:       user.UserId,
		UserPassword: encrypt.HashAndSalt(context.Background(), "123456"),
	}

	if err := orm.DbInstance(context.Background()).Table(orm.USER_AUTH_TABLE).Create(&auth).Error; err != nil {
		log.Error(err)
		return nil
	}

	u.UserInfo.UserId = proto.Int64(user.UserId)
	u.AffiliateInfo.UserId = proto.Int64(user.UserId)

	if u.UserInfo.GetUserRole() == int64(pb.UserRole_ROLE_AFFILIATE) {
		if err := orm.DbInstance(nil).Table(orm.AFFILIATE_DETAILS_TABLE).Create(u.AffiliateInfo).Error; err != nil {
			log.Error(err)
			return nil
		}
	}

	if u.NeedLogin {
		token, err := auth_middleware.CreateToken(context.Background(), u.UserInfo.GetUserId(), u.UserInfo.GetUserRole(), false)
		if err != nil {
			log.Error(err)
			return nil
		}

		saveErr := auth_middleware.CreateAuth(context.Background(), u.UserInfo.GetUserId(), token)
		if saveErr != nil {
			log.Error(err)
			return nil
		}
		u.token = token
		u.Token = &pb.Tokens{
			AccessToken:  proto.String(token.AccessToken),
			RefreshToken: proto.String(token.RefreshToken),
		}
	}
	return u
}

func (u *User) TearDown() {
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v WHERE user_id = %v", orm.USER_TABLE, u.UserInfo.GetUserId())).Error; err != nil {
		log.Error(err)
	}
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v WHERE user_id = %v", orm.USER_AUTH_TABLE, u.UserInfo.GetUserId())).Error; err != nil {
		log.Error(err)
	}
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v WHERE user_id = %v", orm.AFFILIATE_DETAILS_TABLE, u.UserInfo.GetUserId())).Error; err != nil {
		log.Error(err)
	}

	if u.NeedLogin {
		_, err := auth_middleware.DeleteAuth(context.Background(), u.token.AccessUuid)
		if err != nil {
			log.Error(err)
		}

		_, err = auth_middleware.DeleteRefresh(context.Background(), u.token.RefreshUuid)
		if err != nil {
			log.Error(err)
		}
	}
}
