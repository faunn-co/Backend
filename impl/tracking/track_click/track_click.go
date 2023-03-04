package track_click

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
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

	code := t.c.QueryParam("ref")
	if code == "" || code == "null" {
		return nil, nil
	}

	id, err := t.getAffiliateWithCodeUsingCache(code)
	if err != nil {
		return nil, err
	}
	logger.Info(t.ctx, "referral_code: %v referral id: %v", code, id)

	type Referral struct {
		ReferralId        *int64 `gorm:"primary_key"`
		AffiliateId       *int64
		ReferralClickTime *int64
		ReferralStatus    *int64
	}

	r := Referral{
		ReferralId:        nil,
		AffiliateId:       id,
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

func (t *TrackClick) makeCacheKey(code string) string {
	return fmt.Sprintf("ref_code:%v", code)
}

func (t *TrackClick) getAffiliateWithCodeUsingCache(code string) (*int64, *resp.Error) {
	k := t.makeCacheKey(code)
	if val, err := orm.GET(t.c, t.ctx, k, false); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_REDIS)
	} else if val != nil {
		var redisResp *pb.AffiliateDetailsDb
		jsonErr := json.Unmarshal(val, &redisResp)
		if jsonErr != nil {
			logger.Warn(t.ctx, "getAffiliateWithCodeUsingCache | Fail to unmarshal Redis value of key %v : %v, reading from DB", k, jsonErr)
		} else {
			logger.Warn(t.ctx, "getAffiliateWithCodeUsingCache | Successful | Cached %v", k)
			return redisResp.UserId, nil
		}
	}

	var affiliate *pb.AffiliateDetailsDb
	if err := orm.DbInstance(t.ctx).Raw(orm.GetAffiliateByCodeQuery(), code).Scan(&affiliate).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if affiliate.UserId != nil {
		if err := orm.SET(t.ctx, k, affiliate.UserId, 24*time.Hour); err != nil {
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_REDIS)
		}
	}
	return affiliate.UserId, nil
}
