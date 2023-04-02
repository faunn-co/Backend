package update_referral_by_id

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type UpdateReferralById struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *UpdateReferralById {
	g := new(UpdateReferralById)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "UpdateReferralById Initialized")
	return g
}

func (g *UpdateReferralById) UpdateReferralByIdImpl() (*int64, *resp.Error) {
	id := g.c.Param("id")

	if id == "" {
		err := errors.New("invalid id")
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		rDb *pb.ReferralDb
	)

	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralDetailsByIdQuery(), id).Scan(&rDb).Error; err != nil {
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if rDb == nil {
		err := errors.New("id not found")
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	var newStatus int64
	switch rDb.GetReferralStatus() {
	case int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS):
		newStatus = int64(pb.ReferralStatus_REFERRAL_STATUS_CANCELLED)
	case int64(pb.ReferralStatus_REFERRAL_STATUS_CANCELLED):
		//cancelled but without booking cannot be updated
		if rDb.BookingId == nil {
			err := errors.New("no booking exists, not allow to update status to success")
			return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
		}
		newStatus = int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS)
	default:
		err := errors.New("invalid referral status")
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	//update db
	if err := orm.DbInstance(g.ctx).Exec(orm.UpdateReferralStatusByIdQuery(), newStatus, id).Error; err != nil {
		logger.Error(g.ctx, err)
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	logger.Info(g.ctx, "updated referral status to %v", newStatus)
	return proto.Int64(newStatus), nil
}
