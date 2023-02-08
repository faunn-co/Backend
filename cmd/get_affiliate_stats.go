package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_stats"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateStats(c echo.Context) error {
	if curr, prev, err := get_affiliate_stats.New(c).GetAffiliateStatsImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateStatsResponseJSON(c, curr, prev)
	}
}
