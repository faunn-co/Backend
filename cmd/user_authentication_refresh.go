package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/authentication/user_authentication_refresh"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func UserAuthenticationRefresh(c echo.Context) error {
	if t, err := user_authentication_refresh.New(c).UserAuthenticationRefreshImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.UserAuthenticationRefreshResponseJSON(c, t)
	}
}
