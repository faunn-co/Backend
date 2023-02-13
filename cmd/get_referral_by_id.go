package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/get_referral_by_id"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetReferralById(c echo.Context) error {
	if r, err := get_referral_by_id.New(c).GetReferralByIdImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetReferralByIdResponseJSON(c, r)
	}
}
