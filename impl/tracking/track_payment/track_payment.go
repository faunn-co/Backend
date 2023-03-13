package track_payment

import (
	"context"
	"encoding/json"
	"github.com/aaronangxz/AffiliateManager/impl/verification/referral_verification"
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

type TrackPayment struct {
	c   echo.Context
	ctx context.Context
	req *pb.TrackPaymentRequest
}

func New(c echo.Context) *TrackPayment {
	t := new(TrackPayment)
	t.c = c
	t.ctx = logger.NewCtx(t.c)
	logger.Info(t.ctx, "TrackPayment Initialized")
	return t
}

func (t *TrackPayment) TrackPaymentImpl() (*int64, *resp.Error) {
	if err := t.verifyTrackPayment(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	id, err := t.startPaymentTx()
	if err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}
	return id, nil
}

func (t *TrackPayment) startPaymentTx() (*int64, error) {
	tx := orm.DbInstance(t.ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorMsg(t.ctx, "failed to recover: startPaymentTx")
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	marshaledInfo, jErr := json.Marshal(t.req.GetCustomerInfo())
	if jErr != nil {
		logger.Warn(t.ctx, jErr.Error())
	}

	_, citizen, tourist := t.calculateTicket()

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
		BookingStatus:      proto.Int64(int64(pb.BookingStatus_BOOKING_STATUS_PENDING)),
		BookingDay:         t.req.BookingDay,
		BookingSlot:        t.req.BookingSlot,
		TransactionTime:    proto.Int64(time.Now().Unix()),
		PaymentStatus:      proto.Int64(int64(pb.PaymentStatus_PAYMENT_STATUS_PENDING)),
		CitizenTicketCount: t.req.CitizenTicketCount,
		TouristTicketCount: t.req.TouristTicketCount,
		CitizenTicketTotal: citizen,
		TouristTicketTotal: tourist,
		CustomerInfo:       marshaledInfo,
	}

	//decrement ticket
	if err := tx.Exec(orm.DecrementTicketCountQuery(), t.req.GetCitizenTicketCount(), t.req.GetTouristTicketCount(), t.req.GetBookingDay(), t.req.GetBookingSlot()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startPaymentTx:decrement ticket: %v", err.Error())
		tx.Rollback()
		return nil, err
	}

	//Insert into booking_table
	if err := tx.Table(orm.BOOKING_DETAILS_TABLE).Create(&b).Error; err != nil {
		logger.Warn(t.ctx, "Error during startPaymentTx:create booking: %v", err.Error())
		tx.Rollback()
		return nil, err
	}
	//Update referral_table using referral_id
	if err := tx.Exec(orm.UpdateReferralBookingInfoQuery(), pb.ReferralStatus_REFERRAL_STATUS_PENDING, b.BookingId, b.TransactionTime, 0, t.req.GetReferralId()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startPaymentTx:update referral: %v", err.Error())
		tx.Rollback()
		return nil, err
	}
	logger.Info(t.ctx, "committing startPaymentTx")

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return b.BookingId, nil
}

func (t *TrackPayment) verifyTrackPayment() error {
	t.req = new(pb.TrackPaymentRequest)
	if err := t.c.Bind(t.req); err != nil {
		return err
	}
	if err := referral_verification.New(t.c, t.ctx).VerifyReferralId(t.req.GetReferralId()); err != nil {
		return err
	}
	return nil
}

func (t *TrackPayment) calculateTicket() (*int64, *int64, *int64) {
	citizen := t.req.GetCitizenTicketCount() * CitizenTicket
	tourist := t.req.GetTouristTicketCount() * TouristTicket
	total := citizen + tourist
	return proto.Int64(total), proto.Int64(citizen), proto.Int64(tourist)
}
