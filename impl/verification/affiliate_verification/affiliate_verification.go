package affiliate_verification

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
)

type AffiliateVerification struct {
	c   echo.Context
	ctx context.Context
}

func New(c echo.Context, ctx context.Context) *AffiliateVerification {
	a := new(AffiliateVerification)
	a.c = c
	a.ctx = ctx
	return a
}

func (u *AffiliateVerification) VerifyEntityName(name string) error {
	if name == "" {
		return nil
	}
	var user *pb.AffiliateDetailsDb
	if err := orm.DbInstance(u.ctx).Raw(orm.GetAffiliateInfoWithEntityNameQuery(), name).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("affiliate not found")
	}
	return nil
}

func (u *AffiliateVerification) VerifyReferralCode(code string) error {
	if code == "" {
		return nil
	}
	var user *pb.AffiliateDetailsDb
	if err := orm.DbInstance(u.ctx).Raw(orm.GetAffiliateInfoWithReferralCodeQuery(), code).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("affiliate not found")
	}
	return nil
}
