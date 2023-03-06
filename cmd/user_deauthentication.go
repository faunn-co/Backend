package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/authentication/user_deauthentication"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func UserDeAuthentication(c echo.Context) error {
	if err := user_deauthentication.New(c).UserDeAuthenticationImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.UserDeAuthenticationResponseJSON(c)
	}
}
