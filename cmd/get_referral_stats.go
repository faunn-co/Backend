package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/get_referral_stats"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetReferralStats(c echo.Context) error {
	if curr, prev, err := get_referral_stats.New(c).GetReferralStatsImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetReferralStatsResponseJSON(c, curr, prev)
	}
}
