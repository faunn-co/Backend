package rollback_checkout

import (
	"context"
	"encoding/json"
	"github.com/aaronangxz/AffiliateManager/impl/verification/booking_verification"
	"github.com/aaronangxz/AffiliateManager/impl/verification/referral_verification"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"time"
)

type RollbackCheckOut struct {
	c   echo.Context
	ctx context.Context
	req *pb.RollbackCheckOutRequest
}

func New(c echo.Context) *RollbackCheckOut {
	t := new(RollbackCheckOut)
	t.c = c
	t.ctx = logger.NewCtx(t.c)
	logger.Info(t.ctx, "RollbackCheckOut Initialized")
	return t
}

func (t *RollbackCheckOut) RollbackCheckOutImpl() *resp.Error {
	if err := t.verifyRollbackCheckOut(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	if err := t.startRollbackCheckOutTx(); err != nil {
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}
	return nil
}

func (t *RollbackCheckOut) verifyRollbackCheckOut() error {
	t.req = new(pb.RollbackCheckOutRequest)
	if err := t.c.Bind(t.req); err != nil {
		return err
	}
	if err := referral_verification.New(t.c, t.ctx).VerifyReferralId(t.req.GetReferralId()); err != nil {
		return err
	}
	return nil
}

func (t *RollbackCheckOut) startRollbackCheckOutTx() error {
	b, verifyErr := booking_verification.New(t.c, t.ctx).VerifyBookingIdAndGetDetails(t.req.GetBookingId())
	if verifyErr != nil {
		logger.Error(t.ctx, verifyErr)
		return verifyErr
	}

	var c []*pb.CustomerInfo
	if err := json.Unmarshal(b.GetCustomerInfo(), &c); err != nil {
		return err
	}

	tx := orm.DbInstance(t.ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorMsg(t.ctx, "failed to recover: startRollbackCheckOutTx")
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	txTime := time.Now().Unix()
	//increment ticket
	if err := tx.Exec(orm.IncrementTicketCountQuery(), b.GetCitizenTicketCount(), b.GetTouristTicketCount(), b.GetBookingDay(), b.GetBookingSlot()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startRollbackCheckOutTx:decrement ticket: %v", err.Error())
		tx.Rollback()
		return err
	}
	//Update booking_table
	if err := tx.Exec(orm.UpdateBookingPostCheckOutQuery(), int64(pb.BookingStatus_BOOKING_STATUS_FAILED), txTime, int64(pb.PaymentStatus_PAYMENT_STATUS_FAILED), t.req.GetBookingId()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startRollbackCheckOutTx:update booking: %v", err.Error())
		tx.Rollback()
		return err
	}
	//Update referral_table using referral_id
	if err := tx.Exec(orm.UpdateReferralBookingInfoRollbackCheckOutQuery(), pb.ReferralStatus_REFERRAL_STATUS_FAILED, t.req.GetReferralId()).Error; err != nil {
		logger.Warn(t.ctx, "Error during startRollbackCheckOutTx:update referral: %v", err.Error())
		tx.Rollback()
		return err
	}
	logger.Info(t.ctx, "committing startRollbackCheckOutTx")
	if err := tx.Commit().Error; err != nil {
		return err
	}
	if err := referral_verification.New(t.c, t.ctx).PurgeReferralDetailsCache(t.req.GetReferralId()); err != nil {
		return err
	}
	if err := booking_verification.New(t.c, t.ctx).PurgeBookingDetailsCache(t.req.GetBookingId()); err != nil {
		return err
	}
	return nil
}
