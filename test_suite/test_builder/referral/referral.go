package referral

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/booking"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/user"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

type Referral struct {
	ReferralDb *pb.ReferralDb
	HasBooking bool
	Affiliate  *user.User
	Booking    *booking.Booking
}

// New Initializes a new Referral struct with all the required information to have a proper referral record
// This includes creating a new affiliate, booking and referral click.
func New() *Referral {
	r := new(Referral)
	r.ReferralDb = new(pb.ReferralDb)
	r.HasBooking = false
	orm.ENV = "TEST"
	return r
}

func (r *Referral) SetAffiliateId(id int64) *Referral {
	r.ReferralDb.AffiliateId = proto.Int64(id)
	return r
}

func (r *Referral) SetReferralClickTime(clickTime int64) *Referral {
	r.ReferralDb.ReferralClickTime = proto.Int64(clickTime)
	return r
}

func (r *Referral) SetReferralStatus(status int64) *Referral {
	r.ReferralDb.ReferralStatus = proto.Int64(status)
	return r
}

func (r *Referral) SetBookingId(id int64) *Referral {
	r.HasBooking = true
	r.ReferralDb.BookingId = proto.Int64(id)
	return r
}

func (r *Referral) SetBookingTime(bookingTime int64) *Referral {
	r.ReferralDb.BookingTime = proto.Int64(bookingTime)
	return r
}

func (r *Referral) SetReferralCommission(commission int64) *Referral {
	r.ReferralDb.ReferralCommission = proto.Int64(commission)
	return r
}

func (r *Referral) filDefaults() *Referral {
	if r.ReferralDb.AffiliateId == nil {
		r.Affiliate = user.New().Build()
		r.ReferralDb.AffiliateId = r.Affiliate.UserInfo.UserId
	}

	if r.ReferralDb.ReferralClickTime == nil {
		r.ReferralDb.ReferralClickTime = proto.Int64(time.Now().Unix())
	}

	if r.ReferralDb.ReferralStatus == nil {
		r.ReferralDb.ReferralStatus = proto.Int64(int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS))
	}

	if r.ReferralDb.BookingId == nil {
		r.Booking = booking.New().Build()
		r.ReferralDb.BookingId = r.Booking.BookingDetails.BookingId
	}

	if r.ReferralDb.ReferralCommission == nil {
		r.ReferralDb.ReferralCommission = proto.Int64(100)
	}
	return r
}

func (r *Referral) Build() *Referral {
	r.filDefaults()

	type Referral struct {
		ReferralId         *int64 `gorm:"primary_key"`
		AffiliateId        *int64
		ReferralClickTime  *int64
		ReferralStatus     *int64
		BookingId          *int64
		BookingTime        *int64
		ReferralCommission *int64
	}

	ref := Referral{
		ReferralId:         nil,
		AffiliateId:        r.ReferralDb.AffiliateId,
		ReferralClickTime:  r.ReferralDb.ReferralClickTime,
		ReferralStatus:     r.ReferralDb.ReferralStatus,
		BookingId:          r.ReferralDb.BookingId,
		BookingTime:        r.ReferralDb.BookingTime,
		ReferralCommission: r.ReferralDb.ReferralCommission,
	}

	if err := orm.DbInstance(context.Background()).Table(orm.REFERRAL_TABLE).Create(&ref).Error; err != nil {
		log.Error(err)
		return nil
	}

	r.ReferralDb.ReferralId = ref.ReferralId
	return r
}

func (r *Referral) TearDown() {
	if err := orm.DbInstance(context.Background()).Exec(fmt.Sprintf("DELETE FROM %v.%v WHERE referral_id = %v", orm.AFFILIATE_MANAGER_TEST_DB, orm.REFERRAL_TABLE, r.ReferralDb.GetReferralId())).Error; err != nil {
		log.Error(err)
	}

	if r.Affiliate != nil {
		r.Affiliate.TearDown()
	}

	if r.Booking != nil {
		r.Booking.TearDown()
	}
}
