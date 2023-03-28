package referral_verification

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"time"
)

type ReferralVerification struct {
	c                          echo.Context
	ctx                        context.Context
	entityName                 string
	referralId                 string
	referralIdBoundedAffiliate string
}

func New(c echo.Context, ctx context.Context) *ReferralVerification {
	r := new(ReferralVerification)
	r.c = c
	r.ctx = ctx
	r.referralId = "ref_id"
	r.referralIdBoundedAffiliate = "ref_id_bound_affiliate"
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
	if err := orm.SET(r.ctx, k, referral, time.Hour); err != nil {
		logger.Error(r.ctx, err)
		return nil
	}
	return nil
}

func (r *ReferralVerification) PurgeReferralDetailsCache(id int64) error {
	k := fmt.Sprintf("%v:%v", r.referralId, id)
	deleted, err := orm.RedisInstance().Del(context.Background(), k).Result()
	if err != nil {
		logger.Error(r.ctx, err)
		return err
	}
	if deleted == 0 {
		logger.Warn(r.ctx, "purgeReferralIdCache| failed to purge cache | key: %v", k)
	}
	return nil
}

func (r *ReferralVerification) VerifyReferralIdBoundedAffiliate(id int64) error {
	if id == 0 {
		return nil
	}

	k := fmt.Sprintf("%v:%v", r.referralIdBoundedAffiliate, id)
	if val, err := orm.GET(r.c, r.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(r.ctx, "VerifyReferralIdBoundedAffiliate | Successful | Cached %v", k)
		return errors.New("referral click has no affiliate bound")
	}

	var referral *pb.ReferralDb
	if err := orm.DbInstance(r.ctx).Raw(orm.GetReferralClickInfo(), id).Scan(&referral).Error; err != nil {
		return err
	}

	if referral == nil {
		return errors.New("referral click not found")
	}
	if referral.AffiliateId == nil {
		if err := orm.SET(r.ctx, k, referral, time.Hour); err != nil {
			logger.Error(r.ctx, err)
			return nil
		}
		return errors.New("referral click has no affiliate bound")
	}
	return nil
}
