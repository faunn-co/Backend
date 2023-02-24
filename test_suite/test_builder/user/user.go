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
	UserInfo      *pb.User
	AffiliateInfo *pb.AffiliateDetailsDb
}

func New() *User {
	u := new(User)
	u.UserInfo = new(pb.User)
	u.AffiliateInfo = new(pb.AffiliateDetailsDb)
	orm.ENV = "TEST"
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
		u.SetUserName("RandomName")
	}

	if u.UserInfo.UserEmail == nil {
		u.SetUserEmail("random@random.com")
	}

	if u.UserInfo.UserContact == nil {
		u.SetUserContact("+6012345678")
	}

	if u.UserInfo.UserRole == nil {
		u.SetUserRole(int64(pb.UserRole_ROLE_AFFILIATE))
	}

	if u.AffiliateInfo.AffiliateType == nil {
		u.SetAffiliateType(int64(pb.AffiliateType_AFFILIATE_TYPE_ACCOMMODATION))
	}

	if u.AffiliateInfo.UniqueReferralCode == nil {
		u.SetUniqueReferralCode("CODE")
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

	if err := orm.DbInstance(nil).Table(orm.USER_TABLE).Create(&user).Error; err != nil {
		log.Error(err)
		return nil
	}

	u.UserInfo.UserId = proto.Int64(user.UserId)
	u.AffiliateInfo.UserId = proto.Int64(user.UserId)

	if err := orm.DbInstance(nil).Table(orm.AFFILIATE_DETAILS_TABLE).Create(u.AffiliateInfo).Error; err != nil {
		log.Error(err)
		return nil
	}
	return u
}

func (u *User) TearDown() {
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v.%v WHERE user_id = %v", orm.AFFILIATE_MANAGER_TEST_DB, orm.USER_TABLE, u.UserInfo.GetUserId())).Error; err != nil {
		log.Error(err)
	}
	if err := orm.DbInstance(nil).Exec(fmt.Sprintf("DELETE FROM %v.%v WHERE user_id = %v", orm.AFFILIATE_MANAGER_TEST_DB, orm.AFFILIATE_DETAILS_TABLE, u.UserInfo.GetUserId())).Error; err != nil {
		log.Error(err)
	}
}
