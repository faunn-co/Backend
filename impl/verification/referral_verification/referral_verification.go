package referral_verification

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
)

type ReferralVerification struct {
	c          echo.Context
	ctx        context.Context
	entityName string
	referralId string
}

func New(c echo.Context, ctx context.Context) *ReferralVerification {
	r := new(ReferralVerification)
	r.c = c
	r.ctx = ctx
	r.referralId = "ref_id"
	return r
}

func (r *ReferralVerification) VerifyReferralId(id int64) error {
	if id == 0 {
		return nil
	}

	k := fmt.Sprintf("%v:%v", r.referralId, id)
	if val, err := orm.GET(r.c, r.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(r.ctx, "VerifyReferralId | Successful | Cached %v", k)
		return errors.New("referral click already has booking bound")
	}

	var referral *pb.ReferralDb
	if err := orm.DbInstance(r.ctx).Raw(orm.GetReferralClickInfo(), id).Scan(&referral).Error; err != nil {
		return err
	}

	if referral == nil {
		return errors.New("referral click not found")
	}
	if referral.BookingId != nil {
		return errors.New("referral click already has booking bound")
	}
	if err := orm.SET(r.ctx, k, referral, 0); err != nil {
		logger.Error(r.ctx, err)
		return nil
	}
	return nil
}
