package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/update_referral_by_id"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func UpdateReferralById(c echo.Context) error {
	if err := update_referral_by_id.New(c).UpdateReferralByIdImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.UpdateReferralByIdResponseJSON(c)
	}
}
