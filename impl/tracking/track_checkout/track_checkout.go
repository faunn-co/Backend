package track_checkout

import (
	"context"
	"encoding/json"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	CitizenTicket = 8800
	TouristTicket = 9800
)

type TrackCheckOut struct {
	c   echo.Context
	ctx context.Context
	req *pb.TrackCheckOutRequest
}

func New(c echo.Context) *TrackCheckOut {
	t := new(TrackCheckOut)
	t.c = c
	t.ctx = logger.NewCtx(t.c)
	logger.Info(t.ctx, "TrackCheckOut Initialized")
	return t
}

func (t *TrackCheckOut) TrackCheckOutImpl() *resp.Error {
	if err := t.verifyTrackCheckOut(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	if err := t.startCheckOutTx(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}
	return nil
}

func (t *TrackCheckOut) verifyTrackCheckOut() error {
	t.req = new(pb.TrackCheckOutRequest)
	if err := t.c.Bind(t.req); err != nil {
		return err
	}
	return nil
}

func (t *TrackCheckOut) startCheckOutTx() error {
	//TODO
	tx := orm.DbInstance(t.ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorMsg(t.ctx, "failed to recover: startCheckOutTx")
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	marshaledInfo, jErr := json.Marshal(t.req.GetCustomerInfo())
	if jErr != nil {
		logger.Warn(t.ctx, jErr.Error())
	}

	citizen, tourist := t.calculateTicket()

	b := struct {
		BookingId          *int64 `gorm:"primary_key"`
		BookingStatus      *int64
		BookingDay         *string
		BookingSlot        *int64
		TransactionTime    *int64
		PaymentStatus      *int64
		CitizenTicketCount *int64
		TouristTicketCount *int64
		CitizenTicketTotal *int64
		TouristTicketTotal *int64
		CustomerInfo       []byte
	}{
		BookingId:          nil,
		BookingStatus:      proto.Int64(int64(pb.BookingStatus_BOOKING_STATUS_SUCCESS)),
		BookingDay:         t.req.BookingDay,
		BookingSlot:        t.req.BookingSlot,
		TransactionTime:    proto.Int64(time.Now().Unix()),
		PaymentStatus:      proto.Int64(int64(pb.PaymentStatus_PAYMENT_STATUS_SUCCESS)),
		CitizenTicketCount: t.req.CitizenTicketCount,
		TouristTicketCount: t.req.TouristTicketCount,
		CitizenTicketTotal: citizen,
		TouristTicketTotal: tourist,
		CustomerInfo:       marshaledInfo,
	}

	//Insert into booking_table
	if err := tx.Create(&b).Error; err != nil {
		logger.Warn(t.ctx, "Error during startCheckOutTx:create booking: %v", err.Error())
		tx.Rollback()
		return err
	}
	//Update referral_table using referral_id
	//TODO commission calculation
	if err := tx.Exec(orm.UpdateReferralBookingInfoQuery(), pb.ReferralStatus_REFERRAL_STATUS_SUCCESS, b.BookingId, b.TransactionTime, "??", t.req.GetReferralId()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startCheckOutTx:update referral: %v", err.Error())
		tx.Rollback()
		return err
	}
	logger.Info(t.ctx, "committing startCheckOutTx")
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (t *TrackCheckOut) calculateTicket() (*int64, *int64) {
	return proto.Int64(t.req.GetCitizenTicketCount() * CitizenTicket), proto.Int64(t.req.GetTouristTicketCount() * TouristTicket)
}

func (t *TrackCheckOut) calculateCommission() *int64 {
	return nil
}
