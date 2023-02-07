package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_details_by_id"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateDetailsById(c echo.Context) error {
	if m, l, err := get_affiliate_details_by_id.New(c).GetAffiliateDetailsByIdImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateDetailsByIdResponseJSON(c, m, l)
	}
}
