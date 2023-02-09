package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_trend"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateTrend(c echo.Context) error {
	if t, err := get_affiliate_trend.New(c).GetAffiliateTrendImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateTrendResponseJSON(c, t)
	}
}
