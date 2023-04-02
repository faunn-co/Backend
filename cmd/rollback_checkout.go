package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/tracking/rollback_checkout"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func RollbackCheckout(c echo.Context) error {
	if err := rollback_checkout.New(c).RollbackCheckOutImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.RollbackCheckoutResponseJSON(c)
	}
}
