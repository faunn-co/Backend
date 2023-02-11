package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/get_referrals_list"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetReferralsList(c echo.Context) error {
	if l, s, e, err := get_referrals_list.New(c).GetReferralListImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetReferralListResponseJSON(c, l, s, e)
	}
}
