package user

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"google.golang.org/protobuf/proto"
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
	return u
}
