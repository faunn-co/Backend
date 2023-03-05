package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/authentication/user_authentication"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func UserAuthentication(c echo.Context) error {
	if a, err := user_authentication.New(c).UserAuthenticationImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.UserAuthenticationResponseJSON(c, a)
	}
}
