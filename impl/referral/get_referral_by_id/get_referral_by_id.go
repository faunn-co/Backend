package get_referral_by_id

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetReferralById struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *GetReferralById {
	g := new(GetReferralById)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GetReferralById Initialized")
	return g
}

func (g *GetReferralById) GetReferralByIdImpl() (*pb.ReferralDetails, *resp.Error) {
	id := g.c.Param("id")

	if id == "" {
		return nil, resp.BuildError(errors.New("invalid id"), pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		r   = new(pb.ReferralDetails)
		rDb *pb.ReferralDb
		b   *pb.BookingDetailsDb
	)

	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralDetailsByIdQuery(), g.c.Param("id")).Scan(&rDb).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if rDb == nil {
		return nil, resp.BuildError(errors.New("id not found"), pb.GlobalErrorCode_ERROR_FAIL)
	}

	r = &pb.ReferralDetails{
		ReferralId:         rDb.ReferralId,
		AffiliateId:        rDb.AffiliateId,
		ReferralClickTime:  rDb.ReferralClickTime,
		ReferralStatus:     rDb.ReferralStatus,
		BookingId:          nil,
		BookingDetails:     nil,
		ReferralCommission: rDb.ReferralCommission,
	}

	if rDb.BookingId == nil {
		return r, nil
	}

	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralBookingDetailsQuery(), rDb.GetBookingId()).Scan(&b).Error; err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	var c []*pb.CustomerInfo
	if err := json.Unmarshal(b.GetCustomerInfo(), &c); err != nil {
		//return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_JSON_UNMARSHAL)
	}

	r.BookingDetails = &pb.BookingDetails{
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
		CustomerInfo:       c,
	}
	return r, nil
}
