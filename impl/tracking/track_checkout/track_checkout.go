package track_checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/impl/email/send_email"
	"github.com/aaronangxz/AffiliateManager/impl/verification/referral_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

const (
	CitizenTicket        = 8800
	TouristTicket        = 9800
	commissionPercentage = 5
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

func (t *TrackCheckOut) TrackCheckOutImpl() (*pb.BookingDetails, *resp.Error) {
	if err := t.verifyTrackCheckOut(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	d, err := t.startCheckOutTx()
	if err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}
	return d, nil
}

func (t *TrackCheckOut) verifyTrackCheckOut() error {
	t.req = new(pb.TrackCheckOutRequest)
	if err := t.c.Bind(t.req); err != nil {
		return err
	}
	if err := referral_verification.New(t.c, t.ctx).VerifyReferralId(t.req.GetReferralId()); err != nil {
		return err
	}
	return nil
}

func (t *TrackCheckOut) startCheckOutTx() (*pb.BookingDetails, error) {
	tx := orm.DbInstance(t.ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorMsg(t.ctx, "failed to recover: startCheckOutTx")
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
	if err := tx.Table(orm.BOOKING_DETAILS_TABLE).Create(&b).Error; err != nil {
		logger.Warn(t.ctx, "Error during startCheckOutTx:create booking: %v", err.Error())
		tx.Rollback()
		return nil, err
	}
	//Update referral_table using referral_id
	if err := tx.Exec(orm.UpdateReferralBookingInfoQuery(), pb.ReferralStatus_REFERRAL_STATUS_SUCCESS, b.BookingId, b.TransactionTime, t.calculateCommission(), t.req.GetReferralId()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startCheckOutTx:update referral: %v", err.Error())
		tx.Rollback()
		return nil, err
	}
	logger.Info(t.ctx, "committing startCheckOutTx")
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	details := &pb.BookingDetails{
		BookingId:          b.BookingId,
		BookingStatus:      b.BookingStatus,
		BookingDay:         b.BookingDay,
		BookingSlot:        b.BookingSlot,
		TransactionTime:    b.TransactionTime,
		PaymentStatus:      b.PaymentStatus,
		CitizenTicketCount: b.CitizenTicketCount,
		TouristTicketCount: b.TouristTicketCount,
		CitizenTicketTotal: b.CitizenTicketTotal,
		TouristTicketTotal: b.TouristTicketTotal,
		CustomerInfo:       t.req.CustomerInfo,
	}

	t.sendConfirmationEmail(details)
	return details, nil
}

func (t *TrackCheckOut) calculateTicket() (*int64, *int64, *int64) {
	citizen := t.req.GetCitizenTicketCount() * CitizenTicket
	tourist := t.req.GetTouristTicketCount() * TouristTicket
	total := citizen + tourist
	return proto.Int64(total), proto.Int64(citizen), proto.Int64(tourist)
}

func (t *TrackCheckOut) calculateCommission() *int64 {
	if err := referral_verification.New(t.c, t.ctx).VerifyReferralIdBoundedAffiliate(t.req.GetReferralId()); err != nil {
		logger.Info(t.ctx, "anonymous click, no commission calculated")
		return proto.Int64(0)
	}
	total, _, _ := t.calculateTicket()
	commission := *total / 100 * commissionPercentage
	return proto.Int64(commission)
}

func (t *TrackCheckOut) sendConfirmationEmail(details *pb.BookingDetails) {
	var slotMap = map[int64]string{
		0: "Corgi - 10.30am to 12:00pm",
		1: "Corgi - 12.30pm to 02:00pm",
		2: "Dogs - 02.30pm to 04:00pm",
		3: "Dogs - 05.00pm to 06.30pm",
	}
	id := strconv.FormatInt(details.GetBookingId(), 10)
	var ticket string

	if details.CitizenTicketCount != nil {
		ticket += fmt.Sprintf("%v x Citizen ", details.GetCitizenTicketCount())
	}
	if details.TouristTicketCount != nil {
		ticket += fmt.Sprintf("%v x Tourist ", details.GetTouristTicketCount())
	}
	send_email.New(t.c).Send(id, details.GetBookingDay(), slotMap[details.GetBookingSlot()], ticket, details.GetCustomerInfo()[0].GetCustomerEmail())
}
