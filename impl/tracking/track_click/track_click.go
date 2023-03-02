package track_click

import (
	"context"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"time"
)

type TrackClick struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *TrackClick {
	t := new(TrackClick)
	t.c = c
	t.ctx = logger.NewCtx(t.c)
	logger.Info(t.ctx, "TrackClick Initialized")
	return t
}

func (t *TrackClick) TrackClickImpl() (*int64, *resp.Error) {
	if err := t.verifyTrackClick(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	if t.c.QueryParam("ref") == "" || t.c.QueryParam("ref") == "null" {
		return nil, nil
	}

	var affiliate *pb.AffiliateDetailsDb
	if err := orm.DbInstance(t.ctx).Raw(orm.GetAffiliateByCodeQuery(), t.c.QueryParam("ref")).Scan(&affiliate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("referral code not found")
		} else {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
		}
	}
	logger.Info(t.ctx, "referral_code: %v referral id: %v", t.c.QueryParam("ref"), affiliate.GetUserId())

	type Referral struct {
		ReferralId        *int64 `gorm:"primary_key"`
		AffiliateId       *int64
		ReferralClickTime *int64
		ReferralStatus    *int64
	}

	r := Referral{
		ReferralId:        nil,
		AffiliateId:       proto.Int64(affiliate.GetUserId()),
		ReferralClickTime: proto.Int64(time.Now().Unix()),
		ReferralStatus:    proto.Int64(int64(pb.ReferralStatus_REFERRAL_STATUS_PENDING)),
	}
	if err := orm.DbInstance(t.ctx).Table(orm.REFERRAL_TABLE).Create(&r).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return r.ReferralId, nil
}

func (t *TrackClick) verifyTrackClick() error {
	if t.c.QueryParam("ref") == "" || t.c.QueryParam("ref") == "null" {
		log.Warn("no ref id")
	}
	return nil
}
