package booking_verification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
)

type BookingVerification struct {
	c          echo.Context
	ctx        context.Context
	entityName string
	bookingId  string
}

func New(c echo.Context, ctx context.Context) *BookingVerification {
	r := new(BookingVerification)
	r.c = c
	r.ctx = ctx
	r.bookingId = "booking_id"
	return r
}

func (r *BookingVerification) VerifyBookingIdAndGetDetails(id int64) (*pb.BookingDetailsDb, error) {
	if id == 0 {
		return nil, errors.New("invalid booking id")
	}

	k := fmt.Sprintf("%v:%v", r.bookingId, id)
	if val, err := orm.GET(r.c, r.ctx, k, false); err != nil {
		return nil, err
	} else if val != nil {
		logger.Info(r.ctx, "VerifyBookingId | Successful | Cached %v", k)
		var redisResp *pb.BookingDetailsDb
		jsonErr := json.Unmarshal(val, &redisResp)
		if jsonErr != nil {
			logger.Warn(r.ctx, "VerifyBookingId | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
		} else {
			logger.Info(r.ctx, "VerifyBookingId | Successful | Cached %v", k)
			return redisResp, nil
		}
	}

	var booking *pb.BookingDetailsDb
	if err := orm.DbInstance(r.ctx).Raw(orm.GetReferralBookingDetailsQuery(), id).Scan(&booking).Error; err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errors.New("booking id not found")
	}

	if err := orm.SET(r.ctx, k, booking, 0); err != nil {
		logger.Warn(r.ctx, err.Error())
	}
	return booking, nil
}

func (r *BookingVerification) PurgeBookingDetailsCache(id int64) error {
	k := fmt.Sprintf("%v:%v", r.bookingId, id)
	deleted, err := orm.RedisInstance().Del(context.Background(), k).Result()
	if err != nil {
		logger.Error(r.ctx, err)
		return err
	}
	if deleted == 0 {
		logger.Warn(r.ctx, "PurgeBookingDetailsCache| failed to purge cache | key: %v", k)
	}
	return nil
}
