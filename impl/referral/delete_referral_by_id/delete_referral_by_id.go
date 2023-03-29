package delete_referral_by_id

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type DeleteReferralById struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context) *DeleteReferralById {
	g := new(DeleteReferralById)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "DeleteReferralById Initialized")
	return g
}

func (g *DeleteReferralById) DeleteReferralByIdImpl() *resp.Error {
	id := g.c.Param("id")

	if id == "" {
		err := errors.New("invalid id")
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		rDb *pb.ReferralDb
	)

	if err := orm.DbInstance(g.ctx).Raw(orm.GetReferralDetailsByIdQuery(), id).Scan(&rDb).Error; err != nil {
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}

	if rDb == nil {
		err := errors.New("id not found")
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	if rDb.GetReferralStatus() == int64(pb.ReferralStatus_REFERRAL_STATUS_SUCCESS) {
		err := errors.New("successful referral cannot be deleted")
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_FAIL)
	}

	if err := orm.DbInstance(g.ctx).Exec(orm.UpdateReferralStatusByIdQuery(), int64(pb.ReferralStatus_REFERRAL_STATUS_DELETED), id).Error; err != nil {
		logger.Error(g.ctx, err)
		return resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	logger.Info(g.ctx, "deleted referral | id: %v", rDb.GetReferralId())
	return nil
}
