package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/authentication/user_registration"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func UserRegistration(c echo.Context) error {
	if err := user_registration.New(c).UserRegistrationImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.UserRegistrationResponseJSON(c)
	}
}
