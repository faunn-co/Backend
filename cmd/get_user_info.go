package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/user/get_user_info"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetUserInfo(c echo.Context) error {
	if a, u, err := get_user_info.New(c).GetUserInfoImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetUserInfoResponseJSON(c, a, u)
	}
}
