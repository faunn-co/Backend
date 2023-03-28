package update_referral_by_id

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
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

func (g *UpdateReferralById) UpdateReferralByIdImpl() *resp.Error {
	id := g.c.Param("id")

	if id == "" {
		err := errors.New("invalid id")
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		rDb *pb.ReferralDb
	)

	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralDetailsByIdQuery(), g.c.Param("id")).Scan(&rDb).Error; err != nil {
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if rDb == nil {
		err := errors.New("id not found")
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	var newStatus int64
	switch rDb.GetReferralStatus() {
	case int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS):
		newStatus = int64(pb.ReferralStatus_REFERRAL_STATUS_CANCELLED)
	case int64(pb.ReferralStatus_REFERRAL_STATUS_CANCELLED):
		//cancelled but without booking cannot be success
		newStatus = int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS)
	default:
		err := errors.New("invalid referral status")
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	//update db

	return nil
}
