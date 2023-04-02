package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/delete_referral_by_id"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func DeleteReferralById(c echo.Context) error {
	if err := delete_referral_by_id.New(c).DeleteReferralByIdImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.DeleteReferralByIdResponseJSON(c)
	}
}
