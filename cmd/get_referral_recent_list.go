package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/get_referral_recent_list"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetReferralRecentList(c echo.Context) error {
	if l, err := get_referral_recent_list.New(c).GetReferralRecentListImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetReferralRecentListResponseJSON(c, l)
	}
}
