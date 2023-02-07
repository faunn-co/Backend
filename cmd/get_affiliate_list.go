package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_list"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateList(c echo.Context) error {
	if l, err := get_affiliate_list.New(c).GetAffiliateListImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateListResponseJSON(c, l)
	}
}
