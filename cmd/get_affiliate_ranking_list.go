package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/affiliate/get_affiliate_ranking_list"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAffiliateRankingList(c echo.Context) error {
	if curr, err := get_affiliate_ranking_list.New(c).GetAffiliateRankingListImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAffiliateRankingListResponseJSON(c, curr)
	}
}
