package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/referral/get_referral_trend"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetReferralTrend(c echo.Context) error {
	if t, err := get_referral_trend.New(c).GetReferralTrendImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetReferralTrendResponseJSON(c, t)
	}
}
