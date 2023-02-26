package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_info"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateInfo(c echo.Context) error {
	if a, u, err := get_affiliate_info.New(c).GetAffiliateInfoImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateInfoResponseJSON(c, a, u)
	}
}
