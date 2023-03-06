package admin

import (
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/user"
)

type Admin struct {
	User              *user.User
	NumberOfReferrals int
	Referrals         []*referral.Referral
}

// New Initializes a new Admin struct with all the required information to have a proper admin
// This includes creating affiliates, bookings and referral clicks.
func New() *Admin {
	a := new(Admin)
	orm.ENV = "TEST"
	return a
}

func (a *Admin) GenerateReferrals(count int) *Admin {
	a.NumberOfReferrals = count
	return a
}

func (a *Admin) filDefaults() *Admin {
	if a.User == nil {
		a.User = user.New().SetUserRole(int64(pb.UserRole_ROLE_ADMIN)).Build()
	}

	if a.Referrals == nil {
		a.Referrals = make([]*referral.Referral, a.NumberOfReferrals)
		for i := 0; i < a.NumberOfReferrals; i++ {
			a.Referrals[i] = referral.New().Build()
		}
	}
	return a
}

func (a *Admin) Build() *Admin {
	a.filDefaults()
	return a
}

func (a *Admin) TearDown() {
	if a.User != nil {
		a.User.TearDown()
	}

	if a.Referrals != nil {
		for _, r := range a.Referrals {
			r.TearDown()
		}
	}
}
