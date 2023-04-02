package affiliate_verification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"time"
)

type AffiliateVerification struct {
	c            echo.Context
	ctx          context.Context
	affiliateId  string
	entityName   string
	referralCode string
}

func New(c echo.Context, ctx context.Context) *AffiliateVerification {
	a := new(AffiliateVerification)
	a.c = c
	a.ctx = ctx
	a.affiliateId = "affiliate_id"
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

	if err := orm.SET(a.ctx, k, user, time.Hour); err != nil {
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

	if err := orm.SET(a.ctx, k, user, time.Hour); err != nil {
		logger.Error(a.ctx, err)
		return nil
	}
	return nil
}

func (a *AffiliateVerification) VerifyAffiliateId(id int64) (*pb.AffiliateDetailsDb, error) {
	if id == 0 {
		return nil, nil
	}

	k := fmt.Sprintf("%v:%v", a.affiliateId, id)
	if val, err := orm.GET(a.c, a.ctx, k, false); err != nil {
		return nil, err
	} else if val != nil {
		var redisResp *pb.AffiliateDetailsDb
		jsonErr := json.Unmarshal(val, &redisResp)
		if jsonErr != nil {
			logger.Warn(a.ctx, "VerifyAffiliateId | Fail to unmarshal Redis value of key %v : %v, reading from API", k, jsonErr)
		} else {
			logger.Info(a.ctx, "VerifyAffiliateId | Successful | Cached %v | %v", k, redisResp)
			return redisResp, nil
		}
	}

	var user *pb.AffiliateDetailsDb
	if err := orm.DbInstance(a.ctx).Raw(orm.GetAffiliateInfoWithAffiliateId(), id).Scan(&user).Error; err != nil {
		logger.Error(a.ctx, err)
		return nil, err
	}
	if user == nil {
		err := errors.New("affiliate not found")
		logger.Error(a.ctx, err)
		return nil, err
	}

	if err := orm.SET(a.ctx, k, user, time.Hour); err != nil {
		logger.Error(a.ctx, err)
		return nil, nil
	}
	return user, nil
}
