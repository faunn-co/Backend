package user

import (
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

type User struct {
	UserId             *int64
	UserName           *string
	UserEmail          *string
	UserContact        *string
	UserRole           *int64
	CreateTimestamp    *int64
	AffiliateType      *int64
	UniqueReferralCode *string
}

func New() *User {
	u := new(User)
	return u
}

func (u *User) SetUserName(name string) *User {
	u.UserName = proto.String(name)
	return u
}

func (u *User) SetUserEmail(email string) *User {
	u.UserEmail = proto.String(email)
	return u
}

func (u *User) SetUserContact(contact string) *User {
	u.UserContact = proto.String(contact)
	return u
}

func (u *User) SetUserRole(role int64) *User {
	u.UserRole = proto.Int64(role)
	return u
}

func (u *User) SetAffiliateType(affiliateType int64) *User {
	u.AffiliateType = proto.Int64(affiliateType)
	return u
}

func (u *User) SetUniqueReferralCode(code string) *User {
	u.UniqueReferralCode = proto.String(code)
	return u
}

func (u *User) filDefaults() *User {
	if u.UserName == nil {
		u.SetUserName("RandomName")
	}

	if u.UserEmail == nil {
		u.SetUserEmail("random@random.com")
	}

	if u.UserContact == nil {
		u.SetUserContact("+6012345678")
	}

	if u.UserRole == nil {
		u.SetUserRole(int64(pb.UserRole_ROLE_AFFILIATE))
	}

	if u.AffiliateType == nil {
		u.SetAffiliateType(int64(pb.AffiliateType_AFFILIATE_TYPE_ACCOMMODATION))
	}

	if u.UniqueReferralCode == nil {
		u.SetUniqueReferralCode("CODE")
	}
	return u
}

func (u *User) Build() *User {
	u.filDefaults()

	type User struct {
		UserId          *int64
		UserName        *string
		UserEmail       *string
		UserContact     *string
		UserRole        *int64
		CreateTimestamp *int64
	}

	user := User{
		UserId:          nil,
		UserName:        u.UserName,
		UserEmail:       u.UserEmail,
		UserContact:     u.UserContact,
		UserRole:        u.UserRole,
		CreateTimestamp: proto.Int64(time.Now().Unix()),
	}

	if err := orm.DbInstance(nil).Table(orm.USER_TABLE).Create(user).Error; err != nil {
		log.Error(err)
		return u
	}

	u.UserId = user.UserId

	affiliate := &pb.AffiliateDetailsDb{
		UserId:             u.UserId,
		AffiliateType:      u.AffiliateType,
		UniqueReferralCode: u.UniqueReferralCode,
	}

	if err := orm.DbInstance(nil).Table(orm.AFFILIATE_DETAILS_TABLE).Create(affiliate).Error; err != nil {
		log.Error(err)
		return u
	}
	return u
}

func (u *User) TearDown() error {
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v.%v WHERE user_id = %v", orm.AFFILIATE_MANAGER_TEST_DB, orm.USER_TABLE, u.UserId)).Error; err != nil {
		log.Error(err)
		return err
	}
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v.%v WHERE user_id = %v", orm.AFFILIATE_MANAGER_TEST_DB, orm.AFFILIATE_DETAILS_TABLE, u.UserId)).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
