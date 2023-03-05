package affiliate_verification

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
)

type AffiliateVerification struct {
	c            echo.Context
	ctx          context.Context
	entityName   string
	referralCode string
}

func New(c echo.Context, ctx context.Context) *AffiliateVerification {
	a := new(AffiliateVerification)
	a.c = c
	a.ctx = ctx
	a.entityName = "entity_name"
	a.referralCode = "ref_code"
	return a
}

func (a *AffiliateVerification) VerifyEntityName(name string) error {
	if name == "" {
		return nil
	}

	k := fmt.Sprintf("%v:%v", a.entityName, name)
	if val, err := orm.GET(a.c, a.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(a.ctx, "VerifyEntityName | Successful | Cached %v", k)
		return nil
	}

	var user *pb.AffiliateDetailsDb
	if err := orm.DbInstance(a.ctx).Raw(orm.GetAffiliateInfoWithEntityNameQuery(), name).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("affiliate not found")
	}

	if err := orm.SET(a.ctx, k, user, 0); err != nil {
		logger.Error(a.ctx, err)
		return nil
	}
	return nil
}

func (a *AffiliateVerification) VerifyReferralCode(code string) error {
	if code == "" {
		return nil
	}

	k := fmt.Sprintf("%v:%v", a.referralCode, code)
	if val, err := orm.GET(a.c, a.ctx, k, false); err != nil {
		return err
	} else if val != nil {
		logger.Info(a.ctx, "VerifyReferralCode | Successful | Cached %v", k)
		return nil
	}

	var user *pb.AffiliateDetailsDb
	if err := orm.DbInstance(a.ctx).Raw(orm.GetAffiliateInfoWithReferralCodeQuery(), code).Scan(&user).Error; err != nil {
		return err
	}
	if user == nil {
		return errors.New("affiliate not found")
	}

	if err := orm.SET(a.ctx, k, user, 0); err != nil {
		logger.Error(a.ctx, err)
		return nil
	}
	return nil
}
